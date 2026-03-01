export type TimeRangeDTO = {
  from: string;
  to: string;
  replicas: number;
};

export type ScheduleDTO = {
  weekdays: Record<string, TimeRangeDTO[]>;
  dates: Record<string, TimeRangeDTO[]>;
  exceptions: string[];
};

export type ContainerPortDTO = {
  containerPort: number;
  protocol?: string;
};

export type EnvVarDTO = {
  name: string;
  value: string;
};

export type ResourceQuantityDTO = {
  memory?: string;
  cpu?: string;
};

export type ResourcesDTO = {
  requests?: ResourceQuantityDTO;
  limits?: ResourceQuantityDTO;
};

export type HTTPGetActionDTO = {
  path: string;
  port: number;
};

export type ProbeDTO = {
  httpGet?: HTTPGetActionDTO;
  initialDelaySeconds?: number;
  periodSeconds?: number;
};

export type ContainerDTO = {
  name: string;
  image: string;
  ports?: ContainerPortDTO[];
  env?: EnvVarDTO[];
  resources?: ResourcesDTO;
  livenessProbe?: ProbeDTO;
  readinessProbe?: ProbeDTO;
};

export type ApplicationDTO = {
  containers: ContainerDTO[];
};

export type CreateScheduleRequestDTO = {
  schedule: ScheduleDTO;
  application: ApplicationDTO;
};

export type GetScheduleResponseDTO = {
  schedule: ScheduleDTO;
  application: ApplicationDTO;
};

