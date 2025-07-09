(() => {
  // internal/pages/annotate.js
  function AnnotationScript(event) {
    const selection = window.getSelection();
    if (!selection || selection.isCollapsed) return;
    const range = selection.getRangeAt(0);
    console.log(range);
    const start = range.startContainer;
    const startNodeName = start.nodeName.toLowerCase();
    const startSpan = startNodeName === "#text" ? start.parentElement : start;
    const startWord = startNodeName === "#text" ? startSpan.dataset.word : 0;
    const startParagraph = startNodeName === "#text" ? startSpan.parentElement : start;
    const startParagraphIndex = startParagraph.dataset.paragraph;
    const end = range.endContainer;
    const endNodeName = end.nodeName.toLowerCase();
    const endSpan = endNodeName === "#text" ? end.parentElement : end;
    const endWord = endNodeName === "#text" ? endSpan.dataset.word : 0;
    const endParagraph = endNodeName === "#text" ? endSpan.parentElement : end;
    const endParagraphIndex = endParagraph.dataset.paragraph;
    const selectedText = range.toString().trim();
    if (!selectedText) return;
    const info = document.querySelector("#info");
    info.textContent = `Selected "${selectedText}" from paragraph ${startParagraphIndex} word ${startWord} to paragraph ${endParagraphIndex} word ${endWord}`;
    window.start_paragraph.value = startParagraphIndex;
    window.start_word.value = startWord;
    window.end_paragraph.value = endParagraphIndex;
    window.end_word.value = endWord;
    window.selected_text.value = selectedText;
    const popover = document.querySelector("#send_annotation");
    popover.showPopover();
  }

  // src.js
  window.handleAnnotateMouseUp = AnnotationScript;
})();
