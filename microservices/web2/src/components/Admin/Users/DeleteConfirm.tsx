import React from 'react';

interface DeleteConfirmProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => void;
  userName: string;
}

const DeleteConfirm: React.FC<DeleteConfirmProps> = ({ 
  isOpen, 
  onClose, 
  onConfirm, 
  userName 
}) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-[20px]">
      {/* Затемнение фона */}
      <div 
        className="absolute inset-0 bg-black/40 backdrop-blur-[2px]" 
        onClick={onClose}
      ></div>
      
      {/* Модальное окно */}
      <div className="relative w-full max-w-[400px] bg-[var(--bg-card)] rounded-[16px] shadow-2xl border border-[var(--border-color)] p-[24px] animate-in fade-in zoom-in-95 duration-200">
        <div className="flex items-start gap-[16px]">
          
          {/* Иконка предупреждения */}
          <div className="w-[48px] h-[48px] rounded-full bg-[var(--status-danger-bg)] flex items-center justify-center flex-shrink-0">
            <i className="fa-solid fa-triangle-exclamation text-[var(--status-danger-text)] text-[20px]"></i>
          </div>
          
          <div className="flex-1">
            <h3 className="text-[18px] font-semibold text-[var(--text-primary)] mb-[8px]">
              Удалить пользователя?
            </h3>
            <p className="text-[14px] text-[var(--text-secondary)] mb-[24px]">
              Вы действительно хотите удалить{' '}
              <strong className="text-[var(--text-primary)]">{userName}</strong>? 
              Это действие нельзя отменить.
            </p>
            
            {/* Кнопки */}
            <div className="flex gap-[12px]">
              <button
                onClick={onClose}
                className="flex-1 px-[16px] py-[10px] rounded-[10px] border border-[var(--border-color)] text-[var(--text-primary)] hover:bg-[var(--bg-hover)] transition-colors font-medium"
              >
                Отмена
              </button>
              <button
                onClick={onConfirm}
                className="flex-1 px-[16px] py-[10px] rounded-[10px] bg-[var(--status-danger-text)] text-[var(--text-inverse)] hover:opacity-90 transition-opacity font-medium"
              >
                Удалить
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DeleteConfirm;