// export type MachineState = {
//   connection: "connecting" | "connected" | "disconnected";

//   deviceAlive: boolean;
//   deviceReady: boolean;

//   status: string;

//   // сообщение от backend
//   message: string | null;

//   beginState: "idle" | "running" | "error" | "done" | "manual";

//   iteration: number | null;

//   error: string | null;
// };

// type Action =
//   | { type: "WS_OPEN" }
//   | { type: "WS_CLOSE" }
//   | { type: "ACK_BOOT" }
//   | {
//       type: "STATUS";
//       status: string;
//       message?: string;
//     }
//   | {
//       type: "STATE";
//       status: string;
//       message?: string;
//       iteration?: number | null;
//       error?: string | null;
//     }
//   | {
//       type: "FORCE_IDLE";
//       message?: string | null;
//     };
    

// const mapStateToBegin = (
//   status: string,
// ): MachineState["beginState"] => {

//   if (
//     status === "DONE" ||
//     status === "END" ||
//     status === "RETURN_DONE"
//   ) {
//     return "done";
//   }

//   if (status === "ERROR") {
//     return "error";
//   }

//   if (
//     status === "WAIT_READY" ||
//     status === "STAND BY"
//   ) {
//     return "idle";
//   }

//   if (status === "MANUAL_MODE") {
//     return "manual";
//   }

//   return "running";
// };

// export const Reducer = (
//   state: MachineState,
//   action: Action,
// ): MachineState => {

//   switch (action.type) {

//     case "WS_OPEN":
//       return {
//         ...state,
//         connection: "connected",
//       };

//     case "WS_CLOSE":
//       return {
//         ...state,

//         connection: "disconnected",

//         deviceAlive: false,

//         deviceReady: false,

//         beginState: "idle",
//       };

//     case "ACK_BOOT":
//       return {
//         ...state,
//         deviceAlive: true,
//       };

//     case "STATUS":
//       return {
//         ...state,

//         deviceReady: action.status === "READY",

//         message:
//           action.message ??
//           state.message,
//       };

//     case "STATE": {

//       const isError =
//         action.status === "ERROR";

//       return {
//         ...state,

//         status: action.status,

//         message:
//           action.message ??
//           null,

//         iteration:
//           action.iteration ??
//           state.iteration,

//         beginState: isError
//           ? "error"
//           : mapStateToBegin(action.status),

//         error: isError
//           ? (
//               action.error ??
//               action.message ??
//               "Ошибка выполнения"
//             )
//           : null,
//       };
//     }

//     default:
//       return state;
//   }
// };