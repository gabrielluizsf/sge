type FrameRequestCallback = (time: number) => void;

export interface KeyboardEvent extends Event {
  key: string;
}

export interface Animator extends Window {
  requestAnimationFrame(callback: FrameRequestCallback): number;
}