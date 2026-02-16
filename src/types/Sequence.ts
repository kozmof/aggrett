export type Factor = "plus" | "minus";

export type SeqFactor = {
  id: string;
  tag: string;
  time: Date;
  value: number;
  factor: Factor;
};

export type Breakdown = Record<string, { store: number; ids: string[] }>;

export type AccumCore = {
  ids: string[];
  time: Date;
  store: number;
};

export type Accum = {
  ids: string[];
  tags: string[];
  time: Date;
  store: number;
  breakdown: Breakdown;
};
