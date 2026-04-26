export function createRobotCommands(sendMessage: (msg: any) => void) {
  return {
    setReady: () => {
      sendMessage({ type: "SET STATUS READY" });
    },

    sendRaw: (msg: any) => {
      sendMessage(msg);
    },
    
    sendBegin: (
      bunker: number,
      record: any,
      turn: number,
      completed_amount:number,
      required_amount: number,
      extraMode: boolean
    ) => {
      sendMessage(
        {
          type: "BEGIN",
          params: {
            shift: record.shift,
            number: record.number,
            seed: record.seed,
            turn: turn,
            completed_amount: completed_amount,
            required_amount: required_amount,
            bunker: bunker,
            gcode: record.gcode,
            extraMode: extraMode,
          },
        }
      );
    },
    sendStatus: () => {
      sendMessage({ type: "STATUS" });
    },
    sendReturn: () => {
      sendMessage({ type: "RETURN" });
    },
  };
}