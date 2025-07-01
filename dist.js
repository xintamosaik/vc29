(() => {
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __esm = (fn, res) => function __init() {
    return fn && (res = (0, fn[__getOwnPropNames(fn)[0]])(fn = 0)), res;
  };
  var __commonJS = (cb, mod) => function __require() {
    return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
  };

  // src.css
  var init_src = __esm({
    "src.css"() {
    }
  });

  // intel/annotate.js
  function AnnotationScript(event) {
    const selection = window.getSelection();
    if (!selection || selection.isCollapsed) return;
    const range = selection.getRangeAt(0);
    console.log(range);
    const start = range.startContainer;
    const startSpan = start.parentElement;
    const startWord = startSpan.dataset.word;
    const startParagraph = startSpan.parentElement;
    const startParagraphIndex = startParagraph.dataset.paragraph;
    const end = range.endContainer;
    const endSpan = end.parentElement;
    const endParagraph = endSpan.parentElement;
    const endWord = endSpan.dataset.word;
    const endParagraphIndex = endParagraph.dataset.paragraph;
    console.log("Selected text from paragraph", startParagraphIndex, "word", startWord, "to paragraph", endParagraphIndex, "word", endWord);
    const popover = document.querySelector("#send_annotation");
    popover.showPopover();
  }
  var init_annotate = __esm({
    "intel/annotate.js"() {
    }
  });

  // src.js
  var require_src = __commonJS({
    "src.js"() {
      init_src();
      init_annotate();
      window.handleAnnotateMouseUp = AnnotationScript;
    }
  });
  require_src();
})();
