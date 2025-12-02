import type { API } from "./api";

export enum TaskStatus {
  Waiting = "waiting",
  Running = "running",
}

export enum TaskResult {
  Succes = "success",
  Failed = "failed",
}

export interface Task {
  uid: string;
  name: string;
  status: TaskStatus;
  nextRun: Date;
  lastStatus?: TaskResult;
  lastRun?: Date;
  lastMessage?: string;
  lastError?: string;
  interval?: number;
}

export interface TaskHistory {
  id: number;
  name: string;
  result: TaskResult;
  runAt: Date;
  error?: string;
  duration: number;
}

export interface TaskHistoryFilter {
  uid?: string;
  result?: TaskResult;
}

export const convertTask = (task: API.Task): Task => {
  return {
    uid: task.uid,
    name: task.name,
    status: task.status as TaskStatus,
    nextRun: new Date(task.next_run),
    lastStatus: task.last_status ? task.last_status as TaskResult : undefined,
    lastRun: task.last_run ? new Date(task.last_run) : undefined,
    lastMessage: task.last_message,
    lastError: task.last_error,
    interval: task.interval,
  };
}

export const convertTasks = (tasks: API.Task[]): Task[] => {
  return tasks.map(convertTask);
}

export const convertTaskHistory = (history: API.TaskHistory): TaskHistory => {
  return {
    id: history.id,
    name: history.name,
    result: history.result as TaskResult,
    runAt: new Date(history.run_at),
    error: history.error,
    duration: history.duration,
  }
}

export const convertTaskHistories = (histories: API.TaskHistory[]): TaskHistory[] => {
  return histories.map(convertTaskHistory)
}
