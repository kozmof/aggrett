export type Factor = "plus" | "minus";

export type SeqFactor = {
  tag: string;
  time: Date;
  value: number;
  factor: Factor;
};

export type SeqSource = {
  version: string;
  sequence: SeqFactor[];
  baseValue: number;
};

export type Accum = {
  tags: string[];
  time: Date;
  store: number;
  breakdown: Record<string, number>;
};
