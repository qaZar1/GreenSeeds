// src/pages/CameraStream.jsx
import React from "react";

const CameraStream = () => {
  return (
    <div style={{ display: "flex", justifyContent: "center", marginTop: 20 }}>
      <img
        src="/camera/stream"
        alt="Camera Stream"
        style={{
          width: "80%",
          maxWidth: "800px",
          border: "2px solid #ccc",
          borderRadius: "8px",
        }}
      />
    </div>
  );
};

export default CameraStream;
