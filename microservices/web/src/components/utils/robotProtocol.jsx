const EOT = "\x04";

export const encodeMsg = msg => msg + EOT;

export const encodeGcode = (gcode, bunker, shift, number, amount, displayText) => {
  const text = Array.isArray(gcode) ? gcode.join("\n") : gcode;
  const cleaned = text.replace(/\r/g, "");

  const code = `
BEGIN ${shift}/${number}/${amount}
BUNKER ${bunker}
\x02${cleaned}\x03${displayText || ""}
`;

  return encodeMsg(code); // добавляем EOT в конце
};
