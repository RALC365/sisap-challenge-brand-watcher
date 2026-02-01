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
      <div className="flex items-center gap-1">
        <input
          type="date"
          value={startDate}
          onChange={(e) => onStartDateChange(e.target.value)}
          className="input w-[120px] text-sm px-2"
          placeholder="Start"
        />
        <span className="text-text-muted text-xs shrink-0">to</span>
        <input
          type="date"
          value={endDate}
          onChange={(e) => onEndDateChange(e.target.value)}
          className="input w-[120px] text-sm px-2"
          placeholder="End"
        />
      </div>
    </div>
  );
}
