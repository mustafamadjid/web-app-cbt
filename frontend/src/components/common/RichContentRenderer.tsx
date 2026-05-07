import katex from "katex";
import React from "react";

import {
  hasRichContent,
  type RichBlock,
  type RichContent,
  type RichInline,
} from "@/types/Content/RichContent";
import { resolveImageUrl } from "@/helper/MediaUrl/resolveMediaUrl";

type RichContentRendererProps = {
  content?: RichContent | null;
  fallbackText?: string | null;
  className?: string;
  paragraphClassName?: string;
  inlineClassName?: string;
  inline?: boolean;
};

function renderInlineContent(item: RichInline, key: string, inlineClassName: string) {
  if (item.type === "math") {
    const html = katex.renderToString(item.latex, {
      throwOnError: false,
      displayMode: item.display === "block",
      strict: "ignore",
    });

    return (
      <span
        key={key}
        className={item.display === "block" ? "block overflow-x-auto py-1" : inlineClassName}
        dangerouslySetInnerHTML={{ __html: html }}
      />
    );
  }

  if (item.type === "image") {
    const src = resolveImageUrl(item.src);
    if (!src) {
      return null;
    }

    return (
      <img
        key={key}
        src={src}
        alt={item.alt ?? ""}
        className="my-2 inline-block max-h-64 max-w-full rounded border border-slate-200 object-contain align-middle"
        loading="lazy"
      />
    );
  }

  const textNode = <>{item.text}</>;
  const marks = item.marks ?? [];

  if (marks.includes("sup")) {
    return <sup key={key}>{textNode}</sup>;
  }
  if (marks.includes("sub")) {
    return <sub key={key}>{textNode}</sub>;
  }

  return <React.Fragment key={key}>{textNode}</React.Fragment>;
}

function renderBlock(
  block: RichBlock,
  index: number,
  paragraphClassName: string,
  inlineClassName: string,
  inline: boolean,
) {
  const children = block.children.map((child, childIndex) =>
    renderInlineContent(child, `child-${index}-${childIndex}`, inlineClassName),
  );

  if (inline) {
    return (
      <span key={`block-${index}`} className={paragraphClassName}>
        {children}
      </span>
    );
  }

  return (
    <p key={`block-${index}`} className={paragraphClassName}>
      {children}
    </p>
  );
}

const RichContentRenderer: React.FC<RichContentRendererProps> = ({
  content,
  fallbackText,
  className = "",
  paragraphClassName = "whitespace-pre-wrap text-sm leading-relaxed",
  inlineClassName = "inline-block align-middle",
  inline = false,
}) => {
  if (!hasRichContent(content)) {
    if (inline) {
      return <span className={[paragraphClassName, className].join(" ")}>{fallbackText ?? ""}</span>;
    }
    return <p className={[paragraphClassName, className].join(" ")}>{fallbackText ?? ""}</p>;
  }

  const WrapperTag = inline ? "span" : "div";

  return (
    <WrapperTag className={className}>
      {content.blocks.map((block, index) => [
        index > 0 && inline ? <br key={`break-${index}`} /> : null,
        renderBlock(block, index, paragraphClassName, inlineClassName, inline),
      ]).flat()}
    </WrapperTag>
  );
};

export default RichContentRenderer;
