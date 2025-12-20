export function isEmptyParagraph(block: { type: string; data: any }): boolean {
    if (block.type !== 'paragraph') return false;
    const html = String(block?.data?.text ?? '');
    const plain = html
        .replace(/<br\s*\/?>/gi, '\n')
        .replace(/&nbsp;/gi, ' ')
        .replace(/<[^>]*>/g, '')
        .trim();
    return plain.length === 0;
}

type OutputData = { time?: number; version?: string; blocks: any[] };

export function normalizeParagraphEmpties(data: OutputData): OutputData {
    const src = data.blocks ?? [];
    if (src.length === 0) return data;

    let start = 0;
    while (start < src.length && isEmptyParagraph(src[start])) start++;

    let end = src.length - 1;
    while (end >= start && isEmptyParagraph(src[end])) end--;

    const trimmed = src.slice(start, end + 1);

    const out: any[] = [];
    let prevWasEmpty = false;
    for (const b of trimmed) {
        const empty = isEmptyParagraph(b);
        if (empty && prevWasEmpty) {
            continue;
        }
        out.push(b);
        prevWasEmpty = empty;
    }

    return { ...data, blocks: out };
}
