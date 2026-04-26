#!/usr/bin/env python3
"""
Универсальный подсчёт семян на агровате (минеральная вата в лотке).
Исключает область белого пластика; разделяет слипшиеся кластеры через watershed.

Версия с обработкой по тайлам (tiles) для лучшего учета локального освещения/контраста.
"""

from __future__ import annotations

import cv2
import numpy as np
from pathlib import Path
import argparse
from skimage.feature import peak_local_max

def process_with_l_only(bgr: np.ndarray) -> np.ndarray:
    """
    ТОЛЬКО L-КАНАЛ.
    Обрабатывает один фрагмент изображения.
    """
    # Проверка на пустой вход
    if bgr.size == 0:
        return bgr

    lab = cv2.cvtColor(bgr, cv2.COLOR_BGR2LAB)
    L = lab[:, :, 0]
    
    # 1. Усиливаем контраст
    L_enhanced = cv2.convertScaleAbs(L, alpha=2, beta=0)
    L_enhanced = np.clip(L_enhanced, 0, 255)

    # Blackhat помогает выделить темные объекты на светлом фоне (или наоборот, в зависимости от структуры)
    # В данном случае, скорее всего, пытаемся выделить тени или текстуру семян
    L_enhanced = cv2.morphologyEx(L_enhanced, cv2.MORPH_BLACKHAT, cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (5, 5)))

    # Нормализация гистограммы внутри тайла - это ключевой момент при работе с тайлами!
    # Это выравнивает яркость каждого кусочка независимо от соседей.
    L_enhanced = cv2.normalize(L_enhanced, None, 0, 255, cv2.NORM_MINMAX)

    # 2. Лёгкое сглаживание
    L_enhanced = cv2.GaussianBlur(L_enhanced, (3, 3), 0)

    # 3. "Прибиваем белое" (пороговая обработка)
    # Внимание: так как мы сделали normalize, значения теперь 0-255 локально.
    # Порог 50 может работать иначе, чем на полном изображении.
    mask = L_enhanced > 50
    L_enhanced[mask] = 255

    # 4. Дополнительное усиление контраста оставшихся темных областей
    L_enhanced = cv2.convertScaleAbs(L_enhanced, alpha=1.5, beta=0)

    # 5. Утолщение (морфологическое расширение)
    kernel = cv2.getStructuringElement(cv2.MORPH_ELLIPSE, (5, 5))
    L_enhanced = cv2.dilate(L_enhanced, kernel, iterations=1)

    return L_enhanced

def keep_white_make_gray_black(img: np.ndarray, thresh: int = 200) -> np.ndarray:
    """
    Оставляет белое (>= thresh) белым, всё остальное делает чёрным.
    """
    # Проверка
    if img is None or img.size == 0:
        return img

    # Бинаризация
    _, binary = cv2.threshold(img, thresh, 255, cv2.THRESH_BINARY)

    return binary

def process_in_tiles(bgr: np.ndarray, tile_h: int = 200, tile_w: int = 200) -> np.ndarray:
    """
    Разбивает изображение на тайлы размером tile_h x tile_w,
    обрабатывает каждый тайл через process_with_l_only
    и собирает обратно.
    """
    h, w = bgr.shape[:2]
    
    # Создаем пустой массив для результата такого же типа и размера, как входной L-канал (одноканальный)
    # Изначально заполняем нулями или 255, в зависимости от желаемого фона стыков. 
    # Пусть будет 0, так как семена темные на светлом после инверсии логики, но тут у нас L_enhanced где семена?
    # Судя по коду: mask > 50 becomes 255. Значит фон белый, объекты темные? 
    # Нет, смотрим внимательно:
    # Blackhat выделяет темные детали. Normalize растягивает.
    # Если пиксель > 50, он становится 255 (белым).
    # Значит, то что осталось <= 50 - это темные объекты (семена?).
    # Но потом dilate.
    # В итоге мы возвращаем одноканальное изображение.
    
    result = np.zeros((h, w), dtype=np.uint8)
    
    # Проходим по сетке
    for y in range(0, h, tile_h):
        for x in range(0, w, tile_w):
            # Вырезаем тайл
            # Важно не выйти за границы
            y_end = min(y + tile_h, h)
            x_end = min(x + tile_w, w)
            
            tile_bgr = bgr[y:y_end, x:x_end]
            
            # Обрабатываем тайл
            processed_tile = process_with_l_only(tile_bgr)

            # processed_tile = cv2.bitwise_not(processed_tile)

            processed_tile = keep_white_make_gray_black(processed_tile, thresh=200)
            
            # Вставляем обратно
            result[y:y_end, x:x_end] = processed_tile
            
    return result

def count_seeds_watershed_advanced(binary: np.ndarray):
    """
    Улучшенный watershed с поиском локальных максимумов
    Возвращает: (количество, markers)
    """
    binary = (binary > 0).astype(np.uint8) * 255

    # Убираем шум
    kernel = np.ones((3, 3), np.uint8)
    binary = cv2.morphologyEx(binary, cv2.MORPH_OPEN, kernel, iterations=1)

    # Distance transform
    dist = cv2.distanceTransform(binary, cv2.DIST_L2, 5)

    # Нормализация (важно!)
    dist_norm = cv2.normalize(dist, None, 0, 1.0, cv2.NORM_MINMAX)

    # 🔥 ИЩЕМ ЛОКАЛЬНЫЕ МАКСИМУМЫ
    # Делаем dilation и сравниваем
    kernel = np.ones((7, 7), np.uint8)
    local_max = cv2.dilate(dist_norm, kernel)
    peaks = (dist_norm == local_max)

    # Убираем слабые пики
    peaks = peaks & (dist_norm > 0.3)

    peaks = peaks.astype(np.uint8) * 255

    # Маркеры
    num_labels, markers = cv2.connectedComponents(peaks)

    # Watershed
    color = cv2.cvtColor(binary, cv2.COLOR_GRAY2BGR)
    markers = cv2.watershed(color, markers)

    # Считаем
    unique = np.unique(markers)
    seed_count = len(unique[(unique > 1)])  # >1 чтобы убрать фон

    return seed_count, markers

def draw_detected_seeds(binary: np.ndarray) -> np.ndarray:
    """
    Рисует каждое найденное семя отдельным цветом.
    """
    binary = (binary > 0).astype(np.uint8)

    num_labels, labels = cv2.connectedComponents(binary)

    # Создаём цветное изображение
    h, w = labels.shape
    debug_img = np.zeros((h, w, 3), dtype=np.uint8)

    # Генератор случайных цветов
    rng = np.random.default_rng(42)

    for label in range(1, num_labels):  # 0 — фон
        mask = labels == label

        color = rng.integers(0, 255, size=3, dtype=np.uint8)
        debug_img[mask] = color

        # Центр объекта (для номера)
        ys, xs = np.where(mask)
        if len(xs) > 0:
            cx = int(xs.mean())
            cy = int(ys.mean())

            cv2.putText(
                debug_img,
                str(label),
                (cx, cy),
                cv2.FONT_HERSHEY_SIMPLEX,
                0.4,
                (255, 255, 255),
                1,
                cv2.LINE_AA,
            )

    return debug_img

def split_seeds_without_watershed(binary: np.ndarray):
    binary = (binary > 0).astype(np.uint8) * 255

    contours, _ = cv2.findContours(binary, cv2.RETR_EXTERNAL, cv2.CHAIN_APPROX_SIMPLE)

    centers = []

    for cnt in contours:
        area = cv2.contourArea(cnt)

        if area < 20:
            continue

        mask = np.zeros(binary.shape, dtype=np.uint8)
        cv2.drawContours(mask, [cnt], -1, 255, -1)

        dist = cv2.distanceTransform(mask, cv2.DIST_L2, 5)

        # 🔥 ВОТ ОНО
        coordinates = peak_local_max(
            dist,
            min_distance=12,
            threshold_abs=0.3 * dist.max(),
            labels=mask
        )

        if len(coordinates) > 0:
            for y, x in coordinates:
                centers.append((x, y))
        else:
            x, y, w, h = cv2.boundingRect(cnt)
            centers.append((x + w // 2, y + h // 2))

    return centers

def draw_seeds_on_original(bgr: np.ndarray, binary: np.ndarray) -> np.ndarray:
    """
    Рисует найденные семена на оригинальном изображении (bbox + центр).
    """
    # Копия оригинала
    output = bgr.copy()

    # Бинарка → 0/1
    binary = (binary > 0).astype(np.uint8)

    # Поиск компонент
    num_labels, labels = cv2.connectedComponents(binary)

    for label in range(1, num_labels):  # 0 — фон
        mask = labels == label

        ys, xs = np.where(mask)
        if len(xs) == 0:
            continue

        # Bounding box
        x_min, x_max = xs.min(), xs.max()
        y_min, y_max = ys.min(), ys.max()

        # Центр
        cx = int(xs.mean())
        cy = int(ys.mean())

        # Рисуем прямоугольник
        cv2.rectangle(
            output,
            (x_min, y_min),
            (x_max, y_max),
            (0, 255, 0),
            1
        )

        # Маленький квадрат в центре
        size = 3
        cv2.rectangle(
            output,
            (cx - size, cy - size),
            (cx + size, cy + size),
            (0, 0, 255),
            -1
        )

        # Номер (по желанию)
        cv2.putText(
            output,
            str(label),
            (cx + 3, cy - 3),
            cv2.FONT_HERSHEY_SIMPLEX,
            0.3,
            (255, 255, 255),
            1,
            cv2.LINE_AA,
        )

    return output

def draw_centers_boxes(bgr: np.ndarray, centers: list, box_size: int = 8) -> np.ndarray:
    """
    Рисует маленькие квадраты (без заливки) вокруг центров
    """
    output = bgr.copy()

    half = box_size // 2

    for i, (x, y) in enumerate(centers):
        x1 = max(0, x - half)
        y1 = max(0, y - half)
        x2 = min(output.shape[1] - 1, x + half)
        y2 = min(output.shape[0] - 1, y + half)

        # ТОНКИЙ квадрат
        cv2.rectangle(output, (x1, y1), (x2, y2), (0, 255, 0), 1)

    return output

def main():
    parser = argparse.ArgumentParser(description="Обработка изображения семян по тайлам")
    parser.add_argument("image", help="Путь к изображению")
    parser.add_argument("--tile-size", type=int, default=200, help="Размер стороны квадратного тайла (по умолчанию 200)")
    parser.add_argument("--output", "-o", default="./exp/channel.jpg", help="Путь к выходному файлу")
    
    args = parser.parse_args()
    
    input_path = Path(args.image)
    output_path = Path(args.output)
    
    # Создаем директорию вывода, если нет
    output_path.parent.mkdir(parents=True, exist_ok=True)

    bgr = cv2.imread(str(input_path))
    if bgr is None:
        print(f"Не удалось загрузить изображение: {input_path}")
        return
    
    print(f"Обработка изображения {bgr.shape[1]}x{bgr.shape[0]} с тайлами {args.tile_size}x{args.tile_size}...")
    
    # Обрабатываем по тайлам
    result_channel = process_in_tiles(bgr, tile_h=args.tile_size, tile_w=args.tile_size)
    cv2.imwrite(str(output_path.parent / "result_channel.jpg"), result_channel)
    centers = split_seeds_without_watershed(result_channel)
    seed_count = len(centers)

    print(f"Количество семян: {seed_count}")

    debug_img = draw_centers_boxes(bgr, centers)
    
    cv2.imwrite(str(output_path), debug_img)
    print(f"Результат сохранен в {output_path}")

if __name__ == "__main__":
    main()