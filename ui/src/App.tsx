import React, { useMemo, useState } from "react";
import { ScheduleEditor } from "./ScheduleEditor";
import { ActiveScheduleView } from "./ActiveScheduleView";
import { createSchedule, updateSchedule } from "./api";
import type { ApplicationDTO, ScheduleDTO } from "./types";

const INITIAL_SCHEDULE: ScheduleDTO = {
  weekdays: {
    monday: [
      { from: "09:00", to: "11:00", replicas: 2 },
      { from: "11:00", to: "19:00", replicas: 3 }
    ],
    tuesday: [{ from: "08:00", to: "17:00", replicas: 1 }],
    wednesday: [],
    thursday: [],
    friday: [],
    saturday: [],
    sunday: []
  },
  dates: {},
  exceptions: []
};

const INITIAL_APPLICATION: ApplicationDTO = {
  containers: [
    {
      name: "app",
      image: "nginx:1.25",
      ports: [{ containerPort: 80, protocol: "TCP" }],
      env: [],
      resources: {
        requests: { memory: "64Mi", cpu: "20m" },
        limits: { memory: "128Mi", cpu: "20m" }
      },
      livenessProbe: {
        httpGet: { path: "/", port: 80 },
        initialDelaySeconds: 30,
        periodSeconds: 10
      },
      readinessProbe: {
        httpGet: { path: "/", port: 80 },
        initialDelaySeconds: 5,
        periodSeconds: 5
      }
    }
  ]
};

export function App() {
  const [schedule, setSchedule] = useState<ScheduleDTO>(INITIAL_SCHEDULE);
  const [application, setApplication] =
    useState<ApplicationDTO>(INITIAL_APPLICATION);
  const [activeId, setActiveId] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);
  const [lastSavedId, setLastSavedId] = useState<string | null>(null);

  const payload = useMemo(
    () => ({
      schedule,
      application
    }),
    [schedule, application]
  );

  const handleScheduleChange = (
    nextSchedule: ScheduleDTO,
    nextApplication: ApplicationDTO
  ) => {
    setSchedule(nextSchedule);
    setApplication(nextApplication);
  };

  const handleCreate = async () => {
    setSaving(true);
    setSaveError(null);
    try {
      const resp = await createSchedule(payload);
      setActiveId(resp.id);
      setLastSavedId(resp.id);
      window.localStorage.setItem("cron-scaler-active-schedule-id", resp.id);
    } catch (e) {
      setSaveError(e instanceof Error ? e.message : "Failed to create schedule");
    } finally {
      setSaving(false);
    }
  };

  const handleUpdate = async () => {
    if (!activeId) {
      setSaveError("Нет активного ID. Сначала создайте расписание или загрузите его.");
      return;
    }
    setSaving(true);
    setSaveError(null);
    try {
      const resp = await updateSchedule(activeId, payload);
      if (!resp.success) {
        throw new Error("Update failed");
      }
      setLastSavedId(activeId);
    } catch (e) {
      setSaveError(e instanceof Error ? e.message : "Failed to update schedule");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="app-shell">
      <header className="app-header">
        <div className="title-block">
          <h1>Расписание масштабирования</h1>
          <p>Удобный ввод интервалов по дням недели и датам.</p>
        </div>
      </header>

      <main className="layout">
        <section>
          <ScheduleEditor
            schedule={schedule}
            application={application}
            onChange={handleScheduleChange}
          />

          <div className="form-section" style={{ marginTop: 10 }}>
            <div className="row-between">
              <div className="stack-tight">
                <span className="small-label">
                  Сохраните актуальное расписание или создайте новое.
                </span>
                {lastSavedId && (
                  <span className="small-label muted">
                    Последний сохранённый ID:{" "}
                    <code>{lastSavedId}</code>
                  </span>
                )}
              </div>
              <div className="row-actions">
                <button
                  type="button"
                  className="button-secondary"
                  onClick={handleUpdate}
                  disabled={saving}
                >
                  <span className="button-icon">⇄</span>
                  {saving ? "Сохраняем..." : "Обновить активное"}
                </button>
                <button
                  type="button"
                  className="button"
                  onClick={handleCreate}
                  disabled={saving}
                >
                  <span className="button-icon">◎</span>
                  {saving ? "Сохраняем..." : "Создать новое"}
                </button>
              </div>
            </div>
            {saveError && (
              <div className="error-text">
                <span>Ошибка: {saveError}</span>
              </div>
            )}
          </div>
        </section>

        <section>
          <ActiveScheduleView
            activeId={activeId}
            onActiveIdChange={setActiveId}
          />
        </section>
      </main>
      <footer className="app-footer">
        <span>Интерфейс управления расписанием масштабирования.</span>
      </footer>
    </div>
  );
}

