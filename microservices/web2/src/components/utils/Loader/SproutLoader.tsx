import "./SproutLoader.css";

export default function SproutLoader() {
  return (
    <div className="flex items-center justify-center h-full">
      <div className="loader-wrapper">

        <svg viewBox="25 25 50 50" className="circle-container">

          {/* rotating loader circle */}
          <circle cx="50" cy="50" r="20" className="circle-loader"></circle>

        </svg>

        <svg viewBox="0 0 100 100" className="plant">

          {/* stem */}
          <line
            className="stem"
            x1="50"
            y1="80"
            x2="50"
            y2="48"
            stroke="#166534"
            strokeWidth="4"
          />

          {/* bigger leaf */}
          <path
            className="leaf leaf-left"
            d="M50 50 C30 40 25 25 45 30 C50 35 50 45 50 50"
            fill="#22c55e"
          />

          {/* smaller leaf */}
          <path
            className="leaf leaf-right"
            d="M50 50 C70 40 75 25 55 30 C50 35 50 45 50 50"
            fill="#4ade80"
          />

        </svg>

      </div>
    </div>
  );
}