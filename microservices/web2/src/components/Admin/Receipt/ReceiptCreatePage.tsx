import React, { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import toast from "react-hot-toast";
import { api } from "../../../api/apiProvider";
import { usePageHeader } from "../../../context/HeaderContext";

import type { Seed } from "../../../types/seed";
import SproutLoader from "../../utils/Loader/SproutLoader";
import ErrorState from "../../pages/ErrorState";

const ReceiptCreatePage: React.FC = () => {

  const { id } = useParams();
  const navigate = useNavigate();

  usePageHeader(
    id ? "Редактирование рецепта" : "Создание рецепта",
    id ? "Изменение рецепта" : "Добавление нового рецепта"
  );

  const [seeds, setSeeds] = useState<Seed[]>([]);

  const [form, setForm] = useState({
    seed: "",
    description: "",
    gcode: ""
  });

  const [errors, setErrors] = useState<Record<string,string>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(false);

  const loadFormData = async () => {
    setLoading(true);
    setError(false);

    try {
      const requests = [
        api.getList("seeds"),
        ...(id ? [api.getOne("receipts", id)] : [])
      ];

      const results = await Promise.allSettled(requests);

      let hasError = false;

      // seeds всегда первый
      const seedsRes = results[0];

      if (seedsRes.status === "fulfilled") {
        setSeeds(seedsRes.value || []);
      } else {
        toast.error("Не удалось загрузить семена");
        hasError = true;
      }

      if (id) {
        const receiptRes = results[1];

        if (receiptRes.status === "fulfilled") {
          const data = receiptRes.value.data;

          setForm({
            seed: data?.seed ?? "",
            description: data?.description ?? "",
            gcode: data?.gcode ?? ""
          });
        } else {
          toast.error("Ошибка загрузки рецепта");
          hasError = true;
        }
      }

      if (hasError) {
        setError(true);
      }

    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFormData();
  }, [id]);

  const update = (name: string, value: string) => {
    setForm(prev => ({
      ...prev,
      [name]: value
    }));

  };

  const validate = () => {

    const newErrors: Record<string,string> = {};

    if (!form.seed) newErrors.seed = "Обязательное поле";
    if (!form.gcode.trim()) newErrors.gcode = "Обязательное поле";

    setErrors(newErrors);

    return Object.keys(newErrors).length === 0;

  };

  const handleSubmit = async (e: React.FormEvent) => {

    e.preventDefault();

    if (!validate()) return;

    const loading = toast.loading(
      id ? "Сохранение рецепта" : "Создание рецепта"
    );

    try {
      if (id) {
        await api.update("receipts", {
          id: Number(id),
          ...form,
          gcode: form.gcode.trim()
        });
      } else {
        await api.create("receipts", {
          ...form,
          gcode: form.gcode.trim()
        });
      }

      toast.success(
        id ? "Рецепт обновлён" : "Рецепт создан",
        { id: loading }
      );

      navigate("/settings/receipts");
    } catch {
      toast.error("Ошибка сохранения", { id: loading });
    }
  };

  if (loading) return <SproutLoader />
  if (error) return <ErrorState onRetry={loadFormData} />;

  return (

    <div className="w-full space-y-[24px]">
      <div className="bg-[var(--bg-card)] border border-[var(--border-color)] rounded-[16px] shadow-sm">
        <form onSubmit={handleSubmit} className="p-[24px] space-y-[20px]">

          {/* seed */}
          <div>

            <label className="block text-[14px] font-medium mb-[8px] text-[var(--text-primary)]">
              Семена *
            </label>

            <div className="relative">

              <select
								value={form.seed}
								disabled={!!id}
								onChange={e => update("seed", e.target.value)}
								className={`
									w-full px-[14px] py-[10px] pr-[36px] rounded-[10px] border appearance-none
									text-[var(--text-primary)]
									${errors.seed ? "border-red-500" : "border-[var(--border-color)]"}
									${id ? "bg-[var(--bg-disabled)] text-[var(--text-secondary)] cursor-not-allowed opacity-70" : "bg-[var(--bg-page)]"}
								`}
							>

                <option value="">Выберите семена</option>

                {seeds.map(s => (
                  <option key={s.seed} value={s.seed}>
                    {s.seed_ru}
                  </option>
                ))}
              </select>
              <i className="fa-solid fa-chevron-down pointer-events-none absolute right-[12px] top-1/2 -translate-y-1/2 text-[12px] text-[var(--text-secondary)]"/>
            </div>

            {errors.seed && (
              <p className="mt-[4px] text-[12px] text-[var(--status-danger-text)]">
                {errors.seed}
              </p>
            )}
          </div>

          {/* description */}
          <div>
            <label className="block text-[14px] font-medium mb-[8px] text-[var(--text-primary)]">
              Описание
            </label>
            <input
              type="text"
              value={form.description}
              onChange={e => update("description", e.target.value)}
              className="w-full px-[14px] py-[10px] rounded-[10px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)]"
            />
          </div>

          {/* gcode */}
          <div>
            <label className="block text-[14px] font-medium mb-[8px] text-[var(--text-primary)]">
              G-code *
            </label>

            <textarea
              value={form.gcode}
              onChange={e => update("gcode", e.target.value)}
              rows={12}
              className={`w-full px-[14px] py-[10px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)] font-mono resize-none
              ${errors.gcode ? "border-red-500" : "border-[var(--border-color)]"}`}
              placeholder={`G0 X0 Y0
G1 X10 Y10
G1 X20 Y5`
							}
            />

            {errors.gcode && (
              <p className="mt-[4px] text-[12px] text-[var(--status-danger-text)]">
                {errors.gcode}
              </p>
            )}

          </div>

          {/* buttons */}
          <div className="flex gap-[12px] pt-[8px]">

            <button
              type="button"
              onClick={() => navigate(-1)}
              className="flex-1 px-[16px] py-[10px] rounded-[10px] border border-[var(--border-color)] text-[var(--text-primary)]"
            >
              Отмена
            </button>

            <button
              type="submit"
              className="flex-1 px-[16px] py-[10px] rounded-[10px] bg-[var(--color-primary)] text-[var(--text-primary)]"
            >
              {id ? "Сохранить" : "Создать"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );

};

export default ReceiptCreatePage;