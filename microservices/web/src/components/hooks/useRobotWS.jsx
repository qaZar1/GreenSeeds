import { useEffect, useRef, useState } from "react";
import { encodeMsg, encodeGcode } from "../utils/robotProtocol";

const WS_URL = "ws://localhost:8001/ws";

const robotStateTranslate = {
  STAND_BY: "Ожидает нажатия кнопки",
  READY: "Готов",
  BUSY: "В процессе выполнения",
  ERR: "Ошибка",
  END: "Задание выполнено",
  RETURN: "Возврат каретки",
  UNKNOWN: "Неизвестно",
  MANUAL_MODE: "Ручной режим",
};

export function useRobotWS(params = {}) {
  const { record } = params;

  const wsRef = useRef(null);
  const [state, setState] = useState("STAND_BY");
  const [isBoot, setIsBoot] = useState(false);
  const [dots, setDots] = useState("");
  const reconnectTimeout = useRef(null);

  // Анимация точек для BUSY
  useEffect(() => {
    if (state !== "BUSY") {
      setDots("");
      return;
    }
    const interval = setInterval(() => {
      setDots(prev => (prev.length >= 3 ? "" : prev + "."));
    }, 500);
    return () => clearInterval(interval);
  }, [state]);

  // Вызов API при END
  useEffect(() => {
    if (state === "END" && record) {
      fetch("/api/reports/add", {
        method: "POST",
        body: JSON.stringify({
          shift: record.shift,
          number: record.number,
          receipt: record.receipt,
          turn: record.turn,
          success: true,
        }),
      })
      .then(res => res.json())
      .then(data => console.log("END API called:", data))
      .catch(err => console.error(err));
    }
  }, [state, record]);

  // --- Функция подключения к WS с авто-статусом и авто-реконнектом ---
  const connectWS = () => {
    const ws = new WebSocket(WS_URL);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connected");
      // Отправляем STATUS при подключении
      ws.send(encodeMsg("BOOT"));
      ws.send(encodeMsg("STATUS"));
    };

    ws.onmessage = (event) => {
      const msg = event.data;
      console.log("WS message:", msg);

      if (msg.includes("ACK BOOT")) setIsBoot(true);
      else if (msg.includes("READY")) setState("READY");
      else if (msg.includes("BEGIN") || msg.includes("ACK")) setState("BUSY");
      else if (msg.includes("END")) setState("END");
      else if (msg.includes("RETURN")) setState("RETURN");
      else if (msg.includes("STAND BY")) setState("STAND_BY");
      else if (msg.includes("ERR")) setState("ERR");
      else if (msg.includes("Disconnected")) setIsBoot(false);
      else if (msg.includes("MANUAL_MODE")) setState("MANUAL_MODE");
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected. Reconnecting in 2s...");
      reconnectTimeout.current = setTimeout(connectWS, 2000);
    };

    ws.onerror = (err) => {
      console.error("WebSocket error:", err);
      ws.close(); // триггерим onclose для реконнекта
    };
  };

  useEffect(() => {
    connectWS();
    return () => {
      clearTimeout(reconnectTimeout.current);
      wsRef.current?.close();
    };
  }, []);

  const displayState = state === "BUSY" ? robotStateTranslate[state] + dots : robotStateTranslate[state] || state;

  return { 
    rawState: state, 
    displayState, 
    sendCommand: (msg) => wsRef.current?.send(encodeMsg(msg)), 
    sendGcode: (gcode, displayText, bunker, shift, number, amount) => 
      wsRef.current?.send(encodeGcode(gcode, displayText, bunker, shift, number, amount)),
    isBoot,
  };
}
