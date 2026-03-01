import React from "react";
import type {
  ApplicationDTO,
  ScheduleDTO,
  TimeRangeDTO
} from "./types";

const WEEKDAYS = [
  "monday",
  "tuesday",
  "wednesday",
  "thursday",
  "friday",
  "saturday",
  "sunday"
];

export type ScheduleEditorProps = {
  schedule: ScheduleDTO;
  application: ApplicationDTO;
  onChange(schedule: ScheduleDTO, application: ApplicationDTO): void;
};

export function ScheduleEditor({
  schedule,
  application,
  onChange
}: ScheduleEditorProps) {
  const updateSchedule = (patch: Partial<ScheduleDTO>) => {
    onChange(
      { ...schedule, ...patch },
      application
    );
  };

  const updateApplication = (patch: Partial<ApplicationDTO>) => {
    onChange(
      schedule,
      { ...application, ...patch }
    );
  };

  const updateWeekdayRanges = (day: string, ranges: TimeRangeDTO[]) => {
    updateSchedule({
      weekdays: {
        ...schedule.weekdays,
        [day]: ranges
      }
    });
  };

  const addWeekdayRange = (day: string) => {
    const ranges = schedule.weekdays[day] ?? [];
    updateWeekdayRanges(day, [
      ...ranges,
      { from: "09:00", to: "18:00", replicas: 1 }
    ]);
  };

  const updateRange = (
    ranges: TimeRangeDTO[],
    index: number,
    patch: Partial<TimeRangeDTO>
  ) => {
    return ranges.map((r, i) => (i === index ? { ...r, ...patch } : r));
  };

  const removeRange = (ranges: TimeRangeDTO[], index: number) =>
    ranges.filter((_, i) => i !== index);

  return (
    <div className="card">
      <div className="card-inner">
        <div className="card-header">
          <div>
            <div className="card-title">
              <span className="card-title-pill">⏱</span>
              Конструктор расписания
            </div>
            <p className="card-subtitle">
              Настройте интервалы масштабирования для сервиса.
            </p>
          </div>
        </div>

        <div className="form-section">
          <div className="form-section-header">
            <span className="form-section-title">Weekdays</span>
            <span className="small-helper">
              Время в формате HH:MM, 24‑часы.
            </span>
          </div>
          <div className="stack">
            {WEEKDAYS.map((day) => {
              const ranges = schedule.weekdays[day] ?? [];
              return (
                <div key={day} className="weekday-row">
                  <div className="weekday-header">
                    <span className="weekday-name">{day}</span>
                    <div className="row-actions">
                      <span className="small-label">
                        {ranges.length === 0
                          ? "нет интервалов"
                          : `${ranges.length} интервал(ов)`}
                      </span>
                      <button
                        type="button"
                        className="button-ghost"
                        onClick={() => addWeekdayRange(day)}
                      >
                        <span className="button-icon">＋</span>
                        Интервал
                      </button>
                    </div>
                  </div>
                  {ranges.length > 0 && (
                    <div className="time-range-list">
                      {ranges.map((r, idx) => (
                        <div key={idx} className="time-range-row">
                          <div className="field-group">
                            <label className="field-label">
                              From <span>*</span>
                            </label>
                            <input
                              className="input"
                              type="time"
                              value={r.from}
                              onChange={(e) =>
                                updateWeekdayRanges(
                                  day,
                                  updateRange(ranges, idx, {
                                    from: e.target.value
                                  })
                                )
                              }
                            />
                          </div>
                          <div className="field-group">
                            <label className="field-label">
                              To <span>*</span>
                            </label>
                            <input
                              className="input"
                              type="time"
                              value={r.to}
                              onChange={(e) =>
                                updateWeekdayRanges(
                                  day,
                                  updateRange(ranges, idx, {
                                    to: e.target.value
                                  })
                                )
                              }
                            />
                          </div>
                          <div className="field-group">
                            <label className="field-label">
                              Replicas <span>*</span>
                            </label>
                            <input
                              className="input"
                              type="number"
                              min={0}
                              value={r.replicas}
                              onChange={(e) =>
                                updateWeekdayRanges(
                                  day,
                                  updateRange(ranges, idx, {
                                    replicas: Number(e.target.value || 0)
                                  })
                                )
                              }
                            />
                          </div>
                          <div className="row-actions">
                            <button
                              type="button"
                              className="button-ghost"
                              onClick={() =>
                                updateWeekdayRanges(
                                  day,
                                  removeRange(ranges, idx)
                                )
                              }
                            >
                              <span className="button-icon">✕</span>
                            </button>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>

        <div className="form-section">
          <div className="form-section-header">
            <span className="form-section-title">Dates</span>
            <span className="small-helper">
              Точечные даты `YYYY-MM-DD` с интервалами.
            </span>
          </div>

          <DatesEditor schedule={schedule} onChange={updateSchedule} />
        </div>

        <div className="form-section">
          <div className="form-section-header">
            <span className="form-section-title">Exceptions</span>
            <span className="small-helper">
              Даты‑исключения, когда расписание не действует.
            </span>
          </div>
          <ExceptionsEditor schedule={schedule} onChange={updateSchedule} />
        </div>

        <div className="form-section">
          <div className="form-section-header">
            <span className="form-section-title">Application</span>
            <span className="small-helper">
              Описание Deployment (контейнер, порт, ресурсы и т.д.).
            </span>
          </div>
          <ApplicationEditor
            application={application}
            onChange={updateApplication}
          />
        </div>
      </div>
    </div>
  );
}

type SchedulePatchFn = (patch: Partial<ScheduleDTO>) => void;

type DatesEditorProps = {
  schedule: ScheduleDTO;
  onChange: (patch: Partial<ScheduleDTO>) => void;
};

function DatesEditor({ schedule, onChange }: DatesEditorProps) {
  const datesEntries = Object.entries(schedule.dates ?? {});

  const setDateRanges = (date: string, ranges: TimeRangeDTO[]) => {
    onChange({
      dates: {
        ...(schedule.dates ?? {}),
        [date]: ranges
      }
    });
  };

  const addDate = () => {
    const today = new Date();
    const iso = today.toISOString().slice(0, 10);
    setDateRanges(iso, [{ from: "09:00", to: "18:00", replicas: 1 }]);
  };

  const removeDate = (date: string) => {
    const next = { ...(schedule.dates ?? {}) };
    delete next[date];
    onChange({ dates: next });
  };

  return (
    <div className="stack">
      <div className="row-between">
        <span className="small-label">
          {datesEntries.length === 0
            ? "Пока ни одной даты."
            : `${datesEntries.length} дат(ы) настроено.`}
        </span>
        <button type="button" className="button-ghost" onClick={addDate}>
          <span className="button-icon">＋</span>
          Дата
        </button>
      </div>

      <div className="dates-list">
        {datesEntries.map(([date, ranges]) => (
          <div key={date} className="weekday-row">
            <div className="weekday-header">
              <div className="stack-tight">
                <span className="weekday-name">{date}</span>
                <span className="small-label">
                  {ranges.length === 0
                    ? "нет интервалов"
                    : `${ranges.length} интервал(ов)`}
                </span>
              </div>
              <div className="row-actions">
                <button
                  type="button"
                  className="button-ghost"
                  onClick={() => removeDate(date)}
                >
                  <span className="button-icon">✕</span>
                </button>
              </div>
            </div>

            <div className="time-range-list">
              {ranges.map((r, idx) => (
                <div key={idx} className="time-range-row">
                  <div className="field-group">
                    <label className="field-label">
                      From <span>*</span>
                    </label>
                    <input
                      className="input"
                      type="time"
                      value={r.from}
                      onChange={(e) =>
                        setDateRanges(
                          date,
                          updateRange(ranges, idx, {
                            from: e.target.value
                          })
                        )
                      }
                    />
                  </div>
                  <div className="field-group">
                    <label className="field-label">
                      To <span>*</span>
                    </label>
                    <input
                      className="input"
                      type="time"
                      value={r.to}
                      onChange={(e) =>
                        setDateRanges(
                          date,
                          updateRange(ranges, idx, {
                            to: e.target.value
                          })
                        )
                      }
                    />
                  </div>
                  <div className="field-group">
                    <label className="field-label">
                      Replicas <span>*</span>
                    </label>
                    <input
                      className="input"
                      type="number"
                      min={0}
                      value={r.replicas}
                      onChange={(e) =>
                        setDateRanges(
                          date,
                          updateRange(ranges, idx, {
                            replicas: Number(e.target.value || 0)
                          })
                        )
                      }
                    />
                  </div>
                  <div className="row-actions">
                    <button
                      type="button"
                      className="button-ghost"
                      onClick={() =>
                        setDateRanges(date, removeRange(ranges, idx))
                      }
                    >
                      <span className="button-icon">✕</span>
                    </button>
                  </div>
                </div>
              ))}
            </div>

            <div className="row-actions">
              <button
                type="button"
                className="button-ghost"
                onClick={() =>
                  setDateRanges(date, [
                    ...ranges,
                    { from: "09:00", to: "18:00", replicas: 1 }
                  ])
                }
              >
                <span className="button-icon">＋</span>
                Интервал
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

type ExceptionsEditorProps = {
  schedule: ScheduleDTO;
  onChange: (patch: Partial<ScheduleDTO>) => void;
};

function ExceptionsEditor({ schedule, onChange }: ExceptionsEditorProps) {
  const exceptions = schedule.exceptions ?? [];

  const update = (next: string[]) => {
    onChange({ exceptions: next });
  };

  const addException = () => {
    const today = new Date();
    const iso = today.toISOString().slice(0, 10);
    update([...exceptions, iso]);
  };

  const updateException = (index: number, value: string) => {
    update(exceptions.map((e, i) => (i === index ? value : e)));
  };

  const removeException = (index: number) => {
    update(exceptions.filter((_, i) => i !== index));
  };

  return (
    <div className="stack">
      <div className="row-between">
        <span className="small-label">
          {exceptions.length === 0
            ? "Пока нет исключений."
            : `${exceptions.length} дат‑исключений.`}
        </span>
        <button
          type="button"
          className="button-ghost"
          onClick={addException}
        >
          <span className="button-icon">＋</span>
          Дата
        </button>
      </div>

      <div className="exceptions-list">
        {exceptions.map((d, idx) => (
          <div key={idx} className="date-row">
            <div className="field-group">
              <label className="field-label">Дата</label>
              <input
                className="input"
                type="date"
                value={d}
                onChange={(e) => updateException(idx, e.target.value)}
              />
            </div>
            <div className="row-actions">
              <button
                type="button"
                className="button-ghost"
                onClick={() => removeException(idx)}
              >
                <span className="button-icon">✕</span>
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

type ApplicationEditorProps = {
  application: ApplicationDTO;
  onChange(app: ApplicationDTO): void;
};

function ApplicationEditor({ application, onChange }: ApplicationEditorProps) {
  const container = application.containers[0] ?? {
    name: "app",
    image: "nginx:1.25",
    ports: [{ containerPort: 80, protocol: "TCP" as const }],
    env: [],
    resources: {
      requests: { memory: "64Mi", cpu: "20m" },
      limits: { memory: "128Mi", cpu: "20m" }
    }
  };

  const updateContainer = (patch: Partial<typeof container>) => {
    const next = { ...container, ...patch };
    onChange({ containers: [next] });
  };

  const ports = container.ports ?? [];
  const env = container.env ?? [];

  return (
    <div className="stack">
      <div className="grid-2">
        <div className="field-group">
          <label className="field-label">
            Имя контейнера <span>*</span>
          </label>
          <input
            className="input"
            value={container.name}
            onChange={(e) => updateContainer({ name: e.target.value })}
          />
        </div>
        <div className="field-group">
          <label className="field-label">
            Образ <span>*</span>
          </label>
          <input
            className="input"
            value={container.image}
            onChange={(e) => updateContainer({ image: e.target.value })}
            placeholder="nginx:1.25"
          />
        </div>
      </div>

      <div className="divider" />

      <div className="form-section-header">
        <span className="small-label">Порт и протокол</span>
      </div>
      <div className="grid-3">
        <div className="field-group">
          <label className="field-label">
            containerPort <span>*</span>
          </label>
          <input
            className="input"
            type="number"
            min={1}
            value={ports[0]?.containerPort ?? 80}
            onChange={(e) =>
              updateContainer({
                ports: [
                  {
                    containerPort: Number(e.target.value || 0),
                    protocol: ports[0]?.protocol ?? "TCP"
                  }
                ]
              })
            }
          />
        </div>
        <div className="field-group">
          <label className="field-label">Protocol</label>
          <select
            className="select"
            value={ports[0]?.protocol ?? "TCP"}
            onChange={(e) =>
              updateContainer({
                ports: [
                  {
                    containerPort: ports[0]?.containerPort ?? 80,
                    protocol: e.target.value || "TCP"
                  }
                ]
              })
            }
          >
            <option value="TCP">TCP</option>
            <option value="UDP">UDP</option>
          </select>
        </div>
      </div>

      <div className="divider" />

      <div className="form-section-header">
        <span className="small-label">Env‑переменные</span>
      </div>
      <div className="stack">
        {env.length === 0 && (
          <span className="small-label muted">
            Пока нет env. Можно добавить, но не обязательно.
          </span>
        )}
        {env.map((e, idx) => (
          <div key={idx} className="grid-3">
            <div className="field-group">
              <label className="field-label">NAME</label>
              <input
                className="input"
                value={e.name}
                onChange={(ev) =>
                  updateContainer({
                    env: env.map((item, i) =>
                      i === idx ? { ...item, name: ev.target.value } : item
                    )
                  })
                }
              />
            </div>
            <div className="field-group">
              <label className="field-label">VALUE</label>
              <input
                className="input"
                value={e.value}
                onChange={(ev) =>
                  updateContainer({
                    env: env.map((item, i) =>
                      i === idx ? { ...item, value: ev.target.value } : item
                    )
                  })
                }
              />
            </div>
            <div className="row-actions">
              <button
                type="button"
                className="button-ghost"
                onClick={() =>
                  updateContainer({
                    env: env.filter((_, i) => i !== idx)
                  })
                }
              >
                <span className="button-icon">✕</span>
              </button>
            </div>
          </div>
        ))}
        <button
          type="button"
          className="button-ghost"
          onClick={() =>
            updateContainer({
              env: [...env, { name: "ENV_VAR", value: "value" }]
            })
          }
        >
          <span className="button-icon">＋</span>
          Env
        </button>
      </div>

      <div className="divider" />

      <div className="form-section-header">
        <span className="small-label">Ресурсы (K8s quantity)</span>
      </div>
      <div className="grid-2">
        <div className="field-group">
          <label className="field-label">Requests (memory, cpu)</label>
          <div className="grid-2">
            <input
              className="input"
              placeholder="64Mi"
              value={container.resources?.requests?.memory ?? ""}
              onChange={(e) =>
                updateContainer({
                  resources: {
                    ...container.resources,
                    requests: {
                      memory: e.target.value,
                      cpu: container.resources?.requests?.cpu ?? ""
                    }
                  }
                })
              }
            />
            <input
              className="input"
              placeholder="250m"
              value={container.resources?.requests?.cpu ?? ""}
              onChange={(e) =>
                updateContainer({
                  resources: {
                    ...container.resources,
                    requests: {
                      memory: container.resources?.requests?.memory ?? "",
                      cpu: e.target.value
                    }
                  }
                })
              }
            />
          </div>
        </div>
        <div className="field-group">
          <label className="field-label">Limits (memory, cpu)</label>
          <div className="grid-2">
            <input
              className="input"
              placeholder="128Mi"
              value={container.resources?.limits?.memory ?? ""}
              onChange={(e) =>
                updateContainer({
                  resources: {
                    ...container.resources,
                    limits: {
                      memory: e.target.value,
                      cpu: container.resources?.limits?.cpu ?? ""
                    }
                  }
                })
              }
            />
            <input
              className="input"
              placeholder="500m"
              value={container.resources?.limits?.cpu ?? ""}
              onChange={(e) =>
                updateContainer({
                  resources: {
                    ...container.resources,
                    limits: {
                      memory: container.resources?.limits?.memory ?? "",
                      cpu: e.target.value
                    }
                  }
                })
              }
            />
          </div>
        </div>
      </div>
    </div>
  );
}

function updateRange(
  ranges: TimeRangeDTO[],
  index: number,
  patch: Partial<TimeRangeDTO>
): TimeRangeDTO[] {
  return ranges.map((r, i) => (i === index ? { ...r, ...patch } : r));
}

function removeRange(
  ranges: TimeRangeDTO[],
  index: number
): TimeRangeDTO[] {
  return ranges.filter((_, i) => i !== index);
}

