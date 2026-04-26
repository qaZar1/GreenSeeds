import { useNavigate } from "react-router-dom";

const NotFoundPage = () => {
  const navigate = useNavigate();

  return (
    <div className="flex items-center justify-center h-full">
      <div className="text-center">

        {/* 404 */}
        <h1 className="text-[120px] font-bold leading-none bg-gradient-to-r from-gray-700 to-gray-400 bg-clip-text text-transparent">
          404
        </h1>

        {/* текст */}
        <p className="mt-4 text-lg text-gray-500">
          Похоже, такой страницы не существует
        </p>

        {/* кнопки */}
        <div className="mt-6 flex items-center justify-center gap-3">

					<button
						onClick={() => navigate(-1)}
						className="px-[16px] py-[8px] rounded-[8px] border border-[var(--border-color)] hover:bg-[var(--bg-hover)] transition-colors text-[var(--text-primary)]"
					>
						← Назад
					</button>

          <button
            onClick={() => navigate("/dashboard")}
            className="px-4 py-2 rounded-xl bg-[var(--color-primary)] text-white hover:opacity-90 transition"
          >
            В Dashboard
          </button>

        </div>

      </div>
    </div>
  );
};

export default NotFoundPage;