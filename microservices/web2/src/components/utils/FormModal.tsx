import React, { useState, useEffect } from "react";
import type { FormField } from "../../types/form";
import { createPortal } from "react-dom";

interface Props {
  isOpen: boolean;
  title: string;
  onClose: () => void;
  onSubmit: (data: Record<string, any>) => void;
  fields: FormField[];
  initialValues?: Record<string, any>;
}

const FormModal: React.FC<Props> = ({
  isOpen,
  title,
  onClose,
  onSubmit,
  fields,
  initialValues = {},
}) => {

  const [formData, setFormData] = useState<Record<string, any>>(initialValues);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [touched, setTouched] = useState<Record<string, boolean>>({});
  const [submitted, setSubmitted] = useState(false);
  const [showPasswords, setShowPasswords] = useState<Record<string, boolean>>({});

  /* ---------------- reset when modal opens ---------------- */

  useEffect(() => {
    const layout = document.getElementById("app-layout");

    if (!layout) return;

    if (isOpen) {
      layout.classList.add("modal-blur");
    } else {
      layout.classList.remove("modal-blur");
    }
  }, [isOpen]);

  useEffect(() => {
    if (isOpen) {
      setFormData(initialValues);
      setErrors({});
      setTouched({});
      setSubmitted(false);
      setShowPasswords({});
    }
  }, [isOpen, initialValues]);

  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "";
    }

    return () => {
      document.body.style.overflow = "";
    };
  }, [isOpen]);

  /* ---------------- helper ---------------- */

  const resolve = <T,>(value: T | ((values: Record<string, any>) => T)): T => {
    if (typeof value === "function") {
      return (value as any)(formData);
    }
    return value;
  };

  /* ---------------- password toggle ---------------- */

  const togglePassword = (name: string) => {
    setShowPasswords(prev => ({
      ...prev,
      [name]: !prev[name]
    }));
  };

  /* ---------------- validation ---------------- */

  const validateField = (field: FormField, value: any) => {
    const minValue = field.min !== undefined ? resolve(field.min) : undefined;
    const maxValue = field.max !== undefined ? resolve(field.max) : undefined;
    const disabled = field.disabled !== undefined ? resolve(field.disabled) : false;

    if (field.required && !disabled && (value === "" || value === undefined)) {
      return "Обязательное поле";
    }

    if (field.type === "number") {
      if (minValue !== undefined && value < minValue) {
        return `Минимум ${minValue}`;
      }

      if (maxValue !== undefined && value > maxValue) {
        return `Максимум ${maxValue}`;
      }
    }

    return "";
  };

  const validate = () => {
    const newErrors: Record<string, string> = {};

    fields.forEach(field => {
      const value = formData[field.name];
      const error = validateField(field, value);

      if (error) newErrors[field.name] = error;
    });

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  /* ---------------- update field ---------------- */

  const update = (name: string, value: any) => {
    const field = fields.find(f => f.name === name);
    if (!field) return;

    const error = validateField(field, value);

    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    setErrors(prev => ({
      ...prev,
      [name]: error
    }));

    setTouched(prev => ({
      ...prev,
      [name]: true
    }));
  };

  /* ---------------- submit ---------------- */

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitted(true);

    if (!validate()) return;

    onSubmit(formData);
  };

  const hasErrors = Object.keys(errors).some(
    key => errors[key] && (touched[key] || submitted)
  );

  if (!isOpen) return null;

  return createPortal(
    <div className="fixed inset-0 z-50 flex items-center justify-center p-[20px] isolate">

      <div
        className="absolute inset-0 bg-black/40"
      />

      <div className="relative w-full max-w-[480px] bg-[var(--bg-card)] rounded-[16px] shadow-2xl border border-[var(--border-color)]">

        {/* header */}
        <div className="flex items-center justify-between px-[24px] py-[20px] border-b border-[var(--border-color)]">
          <h3 className="text-[18px] font-semibold text-[var(--text-primary)]">
            {title}
          </h3>

          <button
            onClick={onClose}
            className="p-[8px] rounded-[8px] text-[var(--text-secondary)] hover:bg-[var(--bg-hover)]"
          >
            <i className="fa-solid fa-xmark" />
          </button>
        </div>

        {/* form */}
        <form onSubmit={handleSubmit} className="p-[24px] space-y-[20px]">

          {fields.map(field => {
            const error = errors[field.name];
            const value = formData[field.name] ?? "";

            const label = resolve(field.label);
            const disabled = field.disabled !== undefined ? resolve(field.disabled) : false;
            const options = field.options ? resolve(field.options) : undefined;

            const showError = error && (touched[field.name] || submitted);

            return (
              <div key={field.name}>

                <label className="block text-[14px] font-medium mb-[8px] text-[var(--text-primary)]">
                  {label}
                  {field.required && " *"}
                </label>

                {/* text */}
                {field.type === "text" && (
                  <input
                    type="text"
                    value={value}
                    placeholder={field.placeholder}
                    disabled={disabled}
                    onChange={e => update(field.name, e.target.value)}
                    className={`w-full px-[14px] py-[10px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)]
                    ${showError ? "border-red-500" : "border-[var(--border-color)]"}
                    disabled:opacity-60 disabled:cursor-not-allowed`}
                  />
                )}

                {/* number */}
                {field.type === "number" && (
                  <input
                    type="number"
                    value={value}
                    disabled={disabled}
                    onChange={e => {
                      const val = e.target.value === "" ? "" : Number(e.target.value);
                      update(field.name, val);
                    }}
                    className={`w-full px-[14px] py-[10px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)]
                    ${showError ? "border-red-500" : "border-[var(--border-color)]"}
                    disabled:opacity-60 disabled:cursor-not-allowed`}
                  />
                )}

                {/* datetime */}
                {field.type === "datetime" && (
                  <input
                    type="datetime-local"
                    value={value}
                    disabled={disabled}
                    onChange={e => update(field.name, e.target.value)}
                    className={`w-full px-[14px] py-[10px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)]
                    ${showError ? "border-red-500" : "border-[var(--border-color)]"}
                    disabled:opacity-60 disabled:cursor-not-allowed`}
                  />
                )}

                {/* select */}
                {field.type === "select" && (
                  <div className="relative">
                    <select
                      value={value}
                      disabled={disabled}
                      onChange={e => update(field.name, e.target.value)}
                      className={`w-full px-[14px] py-[10px] pr-[36px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)] appearance-none
                      ${showError ? "border-red-500" : "border-[var(--border-color)]"}
                      disabled:opacity-60 disabled:cursor-not-allowed`}
                    >
                      {options?.map(o => (
                        <option key={o.label} value={o.value}>
                          {o.label}
                        </option>
                      ))}
                    </select>

                    <i className="fa-solid fa-chevron-down pointer-events-none absolute right-[12px] top-1/2 -translate-y-1/2 text-[12px] text-[var(--text-secondary)]"/>
                  </div>
                )}

                {/* textarea */}
                {field.type === "textarea" && (
                  <textarea
                    value={value}
                    placeholder={field.placeholder}
                    disabled={disabled}
                    rows={6}
                    onChange={e => update(field.name, e.target.value)}
                    className={`w-full px-[14px] py-[10px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)] resize-none font-mono
                    ${showError ? "border-red-500" : "border-[var(--border-color)]"}
                    disabled:opacity-60 disabled:cursor-not-allowed`}
                  />
                )}

                {/* password */}
                {field.type === "password" && (
                  <div className="relative">
                    <input
                      type={showPasswords[field.name] ? "text" : "password"}
                      value={value}
                      placeholder={field.placeholder}
                      disabled={disabled}
                      onChange={e => update(field.name, e.target.value)}
                      className={`w-full px-[14px] py-[10px] pr-[40px] rounded-[10px] border bg-[var(--bg-page)] text-[var(--text-primary)]
                      ${showError ? "border-red-500" : "border-[var(--border-color)]"}
                      disabled:opacity-60 disabled:cursor-not-allowed`}
                    />

                    <button
                      type="button"
                      onClick={() => togglePassword(field.name)}
                      className="absolute right-[10px] top-1/2 -translate-y-1/2 text-[var(--text-secondary)] hover:text-[var(--text-primary)] transition"
                    >
                      <i className={`fa-solid ${showPasswords[field.name] ? "fa-eye" : "fa-eye-slash"}`} />
                    </button>
                  </div>
                )}

                {showError && (
                  <p className="mt-[4px] text-[12px] text-[var(--status-danger-text)]">
                    {error}
                  </p>
                )}

              </div>
            );
          })}

          {/* buttons */}
          <div className="flex gap-[12px] pt-[8px]">

            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-[16px] py-[10px] rounded-[10px] border border-[var(--border-color)] text-[var(--text-primary)]"
            >
              Отмена
            </button>

            <button
              type="submit"
              disabled={hasErrors}
              className="flex-1 px-[16px] py-[10px] rounded-[10px] bg-[var(--color-primary)] text-[var(--text-primary)] disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Сохранить
            </button>

          </div>

        </form>

      </div>

    </div>,
    document.body
  );
};

export default FormModal;