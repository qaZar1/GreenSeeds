import { useCallback, useEffect, useReducer, useRef } from "react";
import type { TaskRecord } from "../../types/task";

export type WSEvent = "STATE" | "ERROR" | "DONE" | "STOP" | "END";

export type WSStep =
  | "WAIT_READY"
  | "BEGIN"
  | "PHOTO"
  | "CONTROL"
  | "RETURN";

export type Progress = {
  current: number;
  total: number;
  percent: number;
};

export type LogEntry = {
  id: number;
  time: string;
  event: WSEvent;
  step?: WSStep;
  message: string;
};

export type WSResponse = {
  event: WSEvent;
  step?: WSStep;
  message?: string;
  iteration?: number;
  progress?: Progress;
  error?: {
    code: string;
    stage?: WSStep;
    message: string;
  };
};

type MachineState = {
  connection: "connecting" | "connected" | "disconnected";
  step: WSStep | null;
  message: string | null;
  iteration: number | null;
  progress: Progress | null;
  error: {
    code: string;
    stage?: WSStep;
    message: string;
  } | null;
  done: boolean;
  stopped: boolean;
  isRunning: boolean;
  logs: LogEntry[];
};

type Action =
  | { type: "WS_OPEN" }
  | { type: "WS_CLOSE" }
  | {
      type: "STATE";
      step: WSStep | null;
      message: string | null;
      iteration: number | null;
      progress: Progress | null;
    }
  | {
      type: "ERROR";
      step: WSStep | null;
      error: {
        code: string;
        stage?: WSStep;
        message: string;
      };
    }
  | {
      type: "DONE";
      message: string;
    }
  | {
      type: "STOP";
      message: string;
    }
  | {
      type: "RESET";
    }
  | {
      type: "END";
      message: string;
    };

const initialState: MachineState = {
  connection: "connecting",
  step: null,
  message: null,
  iteration: null,
  progress: null,
  error: null,
  done: false,
  stopped: false,
  isRunning: false,
  logs: [],
};

function makeLog(
  event: WSEvent,
  message: string,
  step?: WSStep | null
): LogEntry {
  return {
    id: Date.now() + Math.random(),
    time: new Date().toLocaleTimeString(),
    event,
    step: step ?? undefined,
    message,
  };
}

function reducer(state: MachineState, action: Action): MachineState {
  switch (action.type) {
    case "WS_OPEN":
      return {
        ...state,
        connection: "connected",
      };

    case "WS_CLOSE":
      return {
        ...state,
        connection: "disconnected",
        isRunning: false,
      };

    case "RESET":
      return {
        ...initialState,
        connection: state.connection,
        isRunning: true,
      };

    case "STATE":
      return {
        ...state,
        step: action.step,
        message: action.message,
        iteration: action.iteration,
        progress: action.progress,
        done: false,
        stopped: false,
        error: null,
        isRunning: true,

        logs: action.message
          ? [
              ...state.logs,
              makeLog(
                "STATE",
                action.message,
                action.step,
              ),
            ]
          : state.logs,
      };
    case "ERROR":
      return {
        ...state,
        step: action.step ?? state.step,
        error: {
          code: action.error.code,
          stage: action.error.stage,
          message: action.error.message,
        },
        message: action.error.message,
        isRunning: false,
        logs: [
          ...state.logs,
          makeLog(
            "ERROR",
            action.error.message,
            action.error.stage ?? action.step,
          ),
        ],
      };

    case "DONE":
      return {
        ...state,
        done: true,
        stopped: false,
        isRunning: true,
        message: action.message,

        logs: [
          ...state.logs,
          makeLog(
            "DONE",
            action.message,
          ),
        ],
      };

    case "STOP":
      return {
        ...state,
        stopped: true,
        done: false,
        isRunning: false,
        step: null,
        message: action.message,

        logs: [
          ...state.logs,
          makeLog(
            "STOP",
            action.message,
          ),
        ],
      };

    case "END":
      return {
        ...state,
        isRunning: false,
        message: action.message,

        logs: [
          ...state.logs,
          makeLog(
            "END",
            action.message,
          ),
        ],
      };

    default:
      return state;
  }
}

const WS_URL = "/ws";

export function useRobotWS(token: string | null) {
  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const queueRef = useRef<any[]>([]);
  const connectingRef = useRef(false);

  const [state, dispatch] = useReducer(reducer, initialState);

  const sendMessage = useCallback((data: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(data));
      return;
    }

    queueRef.current.push(data);
  }, []);

  const connect = useCallback(() => {
    if (wsRef.current || connectingRef.current) {
      return;
    }

    connectingRef.current = true;

    const ws = new WebSocket(WS_URL);

    wsRef.current = ws;

    ws.onopen = () => {
      connectingRef.current = false;

      dispatch({
        type: "WS_OPEN",
      });

      if (token) {
        sendMessage({
          type: "AUTH",
          token: "Bearer " + token,
        });
      }

      while (queueRef.current.length) {
        const msg = queueRef.current.shift();
        ws.send(JSON.stringify(msg));
      }
    };

    ws.onclose = () => {
      connectingRef.current = false;
      wsRef.current = null;

      dispatch({
        type: "WS_CLOSE",
      });

      reconnectTimerRef.current = setTimeout(() => {
        connect();
      }, 2000);
    };

    ws.onerror = () => {
      ws.close();
    };

    ws.onmessage = (e) => {
      let msg: WSResponse;

      try {
        msg = JSON.parse(e.data);
      } catch {
        return;
      }

      switch (msg.event) {
        case "STATE":
          dispatch({
            type: "STATE",
            step: msg.step ?? null,
            message: msg.message ?? null,
            iteration: msg.iteration ?? null,
            progress: msg.progress ?? null,
          });
          break;

        case "ERROR":
          const rawMessage = msg.error?.message;
          let safeMessage = "Unknown error";
          
          if (typeof rawMessage === 'string') {
              safeMessage = rawMessage;
          } else if (rawMessage && typeof rawMessage === 'object') {
              safeMessage = (rawMessage as any).text || JSON.stringify(rawMessage);
              if (safeMessage === "{}") safeMessage = "Ошибка без описания";
          }

          dispatch({
            type: "ERROR",
            step: msg.step ?? null,
            error: {
              code: msg.error?.code ?? "INTERNAL",
              stage: msg.error?.stage,
              message: safeMessage,
            },
          });
          break;

        case "DONE":
          dispatch({
            type: "DONE",
            message: msg.message ?? "Done",
          });
          break;

        case "STOP":
          dispatch({
            type: "STOP",
            message: msg.message ?? "Stopped",
          });
          break;
        
        case "END":
          dispatch({
            type: "END",
            message: msg.message ?? "Finished",
          });
          break;
      }
    };
  }, [token, sendMessage]);

  const startPlanting = (record: TaskRecord) => {
    dispatch({
      type: "RESET",
    });

    sendMessage({
      type: "START",
      params: {
        shift: record.shift,
        number: record.number,
        recipe: record.recipe,
        required_amount: record.required_amount,
        bunker: record.bunker,
        gcode: record.gcode,
        extraMode: record.extraMode ?? false,
        seed: record.seed,
      },
    });
  };

  const stopPlanting = () => {
    sendMessage({
      type: "STOP",
    });
  };

  const setReady = () => {
    sendMessage({
      type: "SET STATUS READY",
    });
  };

  useEffect(() => {
    connect();

    return () => {
      if (reconnectTimerRef.current) {
        clearTimeout(reconnectTimerRef.current);
      }

      wsRef.current?.close();
    };
  }, [connect]);

  return {
    sendMessage,
    startPlanting,
    stopPlanting,
    setReady,

    connection: state.connection,
    step: state.step,
    message: state.message,
    iteration: state.iteration,
    progress: state.progress,
    error: state.error,
    done: state.done,
    stopped: state.stopped,
    isRunning: state.isRunning,

    logs: state.logs,

    isConnected: state.connection === "connected",
    hasError: !!state.error,
  };
}