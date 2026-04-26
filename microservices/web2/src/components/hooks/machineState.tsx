type MachineState = {
  connection: "connecting" | "connected" | "disconnected";
  deviceReady: boolean;

  status: string;        // сырой STATE (WAIT_READY, DONE и т.д.)
  beginState: "idle" | "running" | "error" | "done";

  iteration: number | null;

  availableActions: string[] | null;
};

type Event =
  | { type: "WS_OPEN" }
  | { type: "WS_CLOSE" }
  | { type: "AUTH_OK" }
  | { type: "STATE"; status: string; iteration?: number }
  | { type: "BOOT_OK" }
  | { type: "DEVICE_ACK" }
  | { type: "ACTIONS"; actions: string[] | null };

export function Reducer(state: MachineState, event: Event): MachineState {
  switch (event.type) {
    case "WS_OPEN":
      return { ...state, connection: "connected" };

    case "WS_CLOSE":
      return { ...state, connection: "disconnected" };

    case "AUTH_OK":
      return state;

    case "BOOT_OK":
    case "DEVICE_ACK":
      return { ...state, deviceReady: true };

    case "ACTIONS":
      return { ...state, availableActions: event.actions };

    case "STATE": {
      const status = event.status;

      return {
        ...state,
        status,
        iteration: event.iteration ?? state.iteration,
        beginState: mapStateToBegin(status),
      };
    }

    default:
      return state;
  }
}

function mapStateToBegin(status: string): MachineState["beginState"] {
  if (["WAIT_READY", "STAND BY"].includes(status)) return "idle";
  if (["BEGIN_ACK"].includes(status)) return "running";
  if (["DONE"].includes(status)) return "done";
  if (["WAIT_ACTION", "ERROR"].includes(status)) return "error";

  return "running";
}