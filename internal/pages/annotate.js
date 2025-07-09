export default function AnnotationScript(event) {
    const selection = window.getSelection();
    if (!selection || selection.isCollapsed) return;

    const range = selection.getRangeAt(0);
    console.log(range)

    // parse start position: start container
    const start = range.startContainer

    // parse start position: start node name. If it's not a text node the cursor is between elements
    const startNodeName = start.nodeName.toLowerCase();
    
    // parse start position: start span and word and index
    const startSpan = startNodeName === "#text" ? start.parentElement : start;
    const startWord = startNodeName === "#text" ? startSpan.dataset.word : 0;
 
    // parse start position: start paragraph and index
    const startParagraph = startNodeName === "#text" ? startSpan.parentElement : start;
    const startParagraphIndex = startParagraph.dataset.paragraph;

    // parse end position: end container
    const end = range.endContainer

    // parse end position: end node name. If it's not a text node the cursor is between elements
    const endNodeName = end.nodeName.toLowerCase();

    // parse end position: end span and word and index
    const endSpan = endNodeName === "#text" ? end.parentElement : end;
    const endWord = endNodeName === "#text" ? endSpan.dataset.word : 0;

    // parse end position: end paragraph and index
    const endParagraph = endNodeName === "#text" ? endSpan.parentElement : end;
    const endParagraphIndex = endParagraph.dataset.paragraph;

    const selectedText = range.toString().trim();
    if (!selectedText) return;

    // Update the info element with the selected text and positions
    const info = document.querySelector("#info");
    info.textContent = `Selected "${selectedText}" from paragraph ${startParagraphIndex} word ${startWord} to paragraph ${endParagraphIndex} word ${endWord}`;

    // Fill the hidden input fields with the selected text and positions
    // We do it dirty be addressing the ids of elements in global scope
    // Don't do this at home, kids!
    window.start_paragraph.value = startParagraphIndex;
    window.start_word.value = startWord;
  
    window.end_paragraph.value = endParagraphIndex;
    window.end_word.value = endWord;
 
    window.selected_text.value = selectedText;

    const popover = document.querySelector("#send_annotation");

    popover.showPopover();

}
