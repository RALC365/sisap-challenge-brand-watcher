interface DateRangeProps {
  startDate: string;
  endDate: string;
  onStartDateChange: (date: string) => void;
  onEndDateChange: (date: string) => void;
  label?: string;
}

export function DateRange({
  startDate,
  endDate,
  onStartDateChange,
  onEndDateChange,
  label,
}: DateRangeProps) {
  return (
    <div className="w-full">
      {label && <label className="label">{label}</label>}
      <div className="flex items-center gap-2">
        <input
          type="date"
          value={startDate}
          onChange={(e) => onStartDateChange(e.target.value)}
          className="input flex-1"
          placeholder="Start date"
        />
        <span className="text-text-muted">to</span>
        <input
          type="date"
          value={endDate}
          onChange={(e) => onEndDateChange(e.target.value)}
          className="input flex-1"
          placeholder="End date"
        />
      </div>
    </div>
  );
}
