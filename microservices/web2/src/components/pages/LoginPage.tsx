import React, { useState, type FormEvent, type ChangeEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";

const LoginPage: React.FC = () => {
  const navigate = useNavigate();

  const { login } = useAuth();

  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);

    try {
      const data = await login({
        username: name,
        password: password,
      });

      if (data.role === "operator") {
        navigate("/choice");
      } else {
        navigate("/shifts");
      }
    } catch (err) {
      setError("Неверное имя пользователя или пароль");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-[var(--bg-page)] px-4">

      {/* Карточка */}
      <div className="w-full max-w-md bg-[var(--bg-card)] border border-[var(--border-color)] rounded-[14px] shadow-lg p-[32px]">

        {/* Заголовок */}
        <div className="text-center mb-[28px]">
          <div className="inline-flex items-center justify-center w-[56px] h-[56px] rounded-full bg-[var(--status-success-bg)] text-[var(--status-success-text)] text-[24px] mb-[12px]">
            <i className="fa-solid fa-seedling"></i>
          </div>

          <h1 className="text-[20px] font-semibold text-[var(--text-primary)]">
            SeedAdmin
          </h1>

          <p className="text-[13px] text-[var(--text-secondary)] mt-[4px]">
            Вход в систему
          </p>
        </div>

        {/* Форма */}
        <form onSubmit={handleSubmit} className="space-y-[16px]">

          {error && (
            <div className="bg-[var(--status-danger-bg)] text-[var(--status-danger-text)] text-[13px] px-[12px] py-[8px] rounded-[8px] text-center">
              {error}
            </div>
          )}

          {/* Имя пользователя */}
          <div className="relative">
            <i className="fa-solid fa-user absolute left-[12px] top-1/2 -translate-y-1/2 text-[var(--text-secondary)] text-[14px]"></i>

            <input
              type="text"
              value={name}
              onChange={(e: ChangeEvent<HTMLInputElement>) =>
                setName(e.target.value)
              }
              placeholder="Имя пользователя"
              required
              disabled={isLoading}
              className="w-full pl-[36px] pr-[12px] py-[10px] rounded-[8px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)] text-[14px] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
            />
          </div>

          {/* Пароль */}
          <div className="relative">
            <i className="fa-solid fa-lock absolute left-[12px] top-1/2 -translate-y-1/2 text-[var(--text-secondary)] text-[14px]"></i>

            <input
              type="password"
              value={password}
              onChange={(e: ChangeEvent<HTMLInputElement>) =>
                setPassword(e.target.value)
              }
              placeholder="Пароль"
              required
              disabled={isLoading}
              className="w-full pl-[36px] pr-[12px] py-[10px] rounded-[8px] border border-[var(--border-color)] bg-[var(--bg-page)] text-[var(--text-primary)] text-[14px] focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]"
            />
          </div>

          {/* Кнопка */}
          <button
            type="submit"
            disabled={isLoading}
            className="w-full mt-[8px] bg-[var(--color-primary)] hover:bg-[var(--color-primary-hover)] text-[var(--text-inverse)] py-[10px] rounded-[8px] font-medium flex items-center justify-center gap-[8px] transition disabled:opacity-60"
          >
            {isLoading ? (
              <i className="fa-solid fa-circle-notch animate-spin"></i>
            ) : (
              <i className="fa-solid fa-right-to-bracket"></i>
            )}

            {isLoading ? "Вход..." : "Войти"}
          </button>

        </form>
      </div>
    </div>
  );
};

export default LoginPage;