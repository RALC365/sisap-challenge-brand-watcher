import { type SelectHTMLAttributes, forwardRef } from 'react';

interface Option {
  value: string;
  label: string;
}

interface SelectProps extends Omit<SelectHTMLAttributes<HTMLSelectElement>, 'onChange'> {
  label?: string;
  options: Option[];
  error?: string;
  onChange?: (value: string) => void;
}

export const Select = forwardRef<HTMLSelectElement, SelectProps>(
  ({ label, options, error, className = '', id, onChange, ...props }, ref) => {
    const selectId = id || label?.toLowerCase().replace(/\s+/g, '-');

    return (
      <div className="w-full">
        {label && (
          <label htmlFor={selectId} className="label">
            {label}
          </label>
        )}
        <select
          ref={ref}
          id={selectId}
          className={`input ${error ? 'border-error focus:ring-error focus:border-error' : ''} ${className}`}
          onChange={(e) => onChange?.(e.target.value)}
          {...props}
        >
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
        {error && (
          <p className="mt-1 text-sm text-error">{error}</p>
        )}
      </div>
    );
  }
);

Select.displayName = 'Select';

interface MultiSelectProps {
  label?: string;
  options: Option[];
  value: string[];
  onChange: (values: string[]) => void;
  placeholder?: string;
}

export function MultiSelect({ label, options, value, onChange, placeholder = 'Select...' }: MultiSelectProps) {
  const toggleOption = (optionValue: string) => {
    if (value.includes(optionValue)) {
      onChange(value.filter((v) => v !== optionValue));
    } else {
      onChange([...value, optionValue]);
    }
  };

  const selectedLabels = options
    .filter((opt) => value.includes(opt.value))
    .map((opt) => opt.label);

  return (
    <div className="w-full relative">
      {label && <label className="label">{label}</label>}
      <div className="relative">
        <div className="input min-h-[38px] flex flex-wrap gap-1 cursor-pointer">
          {selectedLabels.length > 0 ? (
            selectedLabels.map((label, i) => (
              <span key={i} className="inline-flex items-center px-2 py-0.5 rounded bg-primary/10 text-primary text-xs">
                {label}
              </span>
            ))
          ) : (
            <span className="text-text-muted">{placeholder}</span>
          )}
        </div>
        <div className="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-md shadow-lg z-10 max-h-48 overflow-auto hidden group-focus-within:block peer-focus:block">
          {options.map((option) => (
            <label
              key={option.value}
              className="flex items-center px-3 py-2 hover:bg-gray-50 cursor-pointer"
            >
              <input
                type="checkbox"
                checked={value.includes(option.value)}
                onChange={() => toggleOption(option.value)}
                className="mr-2"
              />
              {option.label}
            </label>
          ))}
        </div>
      </div>
    </div>
  );
}
