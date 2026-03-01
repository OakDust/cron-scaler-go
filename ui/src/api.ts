import type {
  CreateScheduleRequestDTO,
  GetScheduleResponseDTO
} from "./types";

const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

async function handleResponse<T>(res: Response): Promise<T> {
  if (!res.ok) {
    let message = `HTTP ${res.status}`;
    try {
      const data = await res.json();
      if (data && typeof data.error === "string") {
        message = data.error;
      }
    } catch {
      // ignore
    }
    throw new Error(message);
  }
  return (await res.json()) as T;
}

export async function createSchedule(
  payload: CreateScheduleRequestDTO
): Promise<{ id: string }> {
  const res = await fetch(`${API_BASE_URL}/v1/schedules`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });
  return handleResponse<{ id: string }>(res);
}

export async function updateSchedule(
  id: string,
  payload: CreateScheduleRequestDTO
): Promise<{ success: boolean }> {
  const res = await fetch(`${API_BASE_URL}/v1/schedules/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json"
    },
    body: JSON.stringify(payload)
  });
  return handleResponse<{ success: boolean }>(res);
}

export async function getSchedule(
  id: string
): Promise<GetScheduleResponseDTO> {
  const res = await fetch(`${API_BASE_URL}/v1/schedules/${id}`);
  return handleResponse<GetScheduleResponseDTO>(res);
}

