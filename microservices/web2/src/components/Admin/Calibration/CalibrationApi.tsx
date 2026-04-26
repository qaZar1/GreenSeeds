import { api } from "../../../api/apiProvider";

export const calibrationApi = {

  handshake: () =>
    api.create("calibration", {}),

  photo: (number: number, sessionId: string) =>
    api.create(
      "takePhoto",
      { numberPhoto: number },
      {
        headers: {
          "X-Calibration-Session": sessionId
        }
      }
    ),

  calculate: (steps: number, sessionId: string) =>
    api.create(
      "calculate",
      { steps: steps },
      {
        headers: {
          "X-Calibration-Session": sessionId
        }
      }
    ),

  save: (sessionId: string) =>
    api.create(
      "saveCalibration",
      {},
      {
        headers: {
          "X-Calibration-Session": sessionId
        }
      }
    )

};