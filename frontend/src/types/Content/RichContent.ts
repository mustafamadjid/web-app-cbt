export type RichTextMark = "sup" | "sub";

export type RichInline =
  | {
      type: "text";
      text: string;
      marks?: RichTextMark[];
    }
  | {
      type: "math";
      latex: string;
      display?: "inline" | "block";
    };

export type RichBlock = {
  type: "paragraph";
  children: RichInline[];
};

export type RichContent = {
  version: number;
  blocks: RichBlock[];
};

export function hasRichContent(value: RichContent | null | undefined): value is RichContent {
  return Boolean(value && Array.isArray(value.blocks) && value.blocks.length > 0);
}

