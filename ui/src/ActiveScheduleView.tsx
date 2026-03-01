import React, { useEffect, useState } from "react";
import { getSchedule } from "./api";
import type { GetScheduleResponseDTO } from "./types";

const ACTIVE_ID_STORAGE_KEY = "cron-scaler-active-schedule-id";

export type ActiveScheduleViewProps = {
  activeId: string | null;
  onActiveIdChange(id: string | null): void;
};

export function ActiveScheduleView({
  activeId,
  onActiveIdChange
}: ActiveScheduleViewProps) {
  const [inputId, setInputId] = useState(activeId ?? "");
  const [data, setData] = useState<GetScheduleResponseDTO | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    setInputId(activeId ?? "");
  }, [activeId]);

  useEffect(() => {
    const stored = window.localStorage.getItem(ACTIVE_ID_STORAGE_KEY);
    if (!activeId && stored) {
      onActiveIdChange(stored);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const loadSchedule = async (id: string) => {
    if (!id) {
      setError("Укажите ID расписания (UUID).");
      return;
    }
    setLoading(true);
    setError(null);
    try {
      const resp = await getSchedule(id);
      setData(resp);
      onActiveIdChange(id);
      window.localStorage.setItem(ACTIVE_ID_STORAGE_KEY, id);
    } catch (e) {
      setData(null);
      setError(e instanceof Error ? e.message : "Failed to load schedule");
    } finally {
      setLoading(false);
    }
  };

  const handleLoadClick = () => {
    void loadSchedule(inputId.trim());
  };

  return (
    <div className="card">
      <div className="card-inner">
        <div className="card-header">
          <div>
            <div className="card-title">
              <span className="card-title-pill">◎</span>
              Активное расписание
            </div>
            <p className="card-subtitle">
              Загрузите сохранённое расписание по его идентификатору.
            </p>
          </div>
          <div className="stack-tight" style={{ alignItems: "flex-end" }}>
            {data ? (
              <div className="status-pill">
                <span className="status-dot" />
                <span>Загружено</span>
              </div>
            ) : (
              <span className="small-label muted">
                Пока ничего не загружено.
              </span>
            )}
          </div>
        </div>

        <div className="form-section">
          <div className="form-section-header">
            <span className="form-section-title">Schedule ID</span>
            <span className="small-helper">
              Идентификатор расписания, полученный при сохранении.
            </span>
          </div>
          <div className="grid-2">
            <input
              className="input"
              placeholder="например, 7b0b0e9b-..."
              value={inputId}
              onChange={(e) => setInputId(e.target.value)}
            />
            <button
              type="button"
              className="button"
              onClick={handleLoadClick}
              disabled={loading}
            >
              <span className="button-icon">↻</span>
              {loading ? "Загружаем..." : "Загрузить"}
            </button>
          </div>
          {error && (
            <div className="error-text">
              <span>Ошибка: {error}</span>
            </div>
          )}
        </div>

        <div className="form-section">
          <div className="form-section-header">
            <span className="form-section-title">JSON</span>
            <span className="small-helper">Текущее содержимое выбранного расписания.</span>
          </div>
          <div className="json-view">
            {data ? (
              <pre>
                {JSON.stringify(
                  {
                    schedule: data.schedule,
                    application: data.application
                  },
                  null,
                  2
                )}
              </pre>
            ) : (
              <pre className="muted">
                {`{\n  "schedule": {\n    \"weekdays\": { ... },\n    \"dates\": { ... },\n    \"exceptions\": []\n  },\n  \"application\": {\n    \"containers\": []\n  }\n}`}
              </pre>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

