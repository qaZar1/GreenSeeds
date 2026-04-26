export type RobotState =
  | "STAND BY"
  | "READY"
  | "BUSY"
  | "END"
  | "RETURN"
  | "WAIT"
  | "ERR"
  | "MANUAL_MODE"
  | "UNKNOWN";

/* ===== Incoming messages ===== */

export type WSErrorMessage = {
  type: "ERR";
  error: string;
  params?: {
    turn?: number;
    bunker?: number;
  };
};

export type WSBootMessage = {
  type: "BOOT";
  status: string;
};

export type WSStatusMessage = {
  type: "STATUS";
  status: RobotState;
};

export type WSBeginMessage = {
  type: "BEGIN";
  status: string;
  payload?: {
    control?: boolean;
  };
  params?: {
    turn?: number;
  };
  error?: string;
};

export type WSDecisionMessage = {
  type: "NEED_DECISION";
  payload: {
    reason: string;
    photo?: string;
  };
};

export type WSSetReadyMessage = {
  type: "SETSTATUS READY";
  status: string;
};

export type WSBunkersMessage = {
  type: "BUNKERS_UPDATE";
  bunkers: {
    bunker: number;
    amount: number;
  }[];
};

export type WSIncomingMessage =
  | WSErrorMessage
  | WSBootMessage
  | WSStatusMessage
  | WSBeginMessage
  | WSDecisionMessage
  | WSSetReadyMessage
  | WSBunkersMessage;

/* ===== Outgoing ===== */

export type WSOutgoingMessage =
  | { type: "BOOT" }
  | { type: "STATUS" }
  | {
      type: "BEGIN";
      params: {
        shift: string;
        number: number;
        seed: string;
        turn: number;
        completed_amount: number;
        required_amount: number;
        bunker: number;
        gcode?: string;
        extraMode?: boolean;
      };
    }
  | {
      type: "DECISION";
      status: "OK" | "NOK";
      solution: string;
    };

/* ===== Guards ===== */

export const isWSMessage = (msg: unknown): msg is WSIncomingMessage => {
  return (
    typeof msg === "object" &&
    msg !== null &&
    "type" in msg &&
    typeof (msg as any).type === "string"
  );
};