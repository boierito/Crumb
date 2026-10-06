# Review of the four supplied screenshots

Scope: user-supplied iPhone screenshots from the current conversation and the
corresponding native SwiftUI source, compared with original WaffleStore/PartyUI
structure. This is a screenshot-based review, not a live flow/accessibility audit.
The supplied images were viewed in the conversation; they contain account data
and are not copied into the repository or redistributed. No simulator rendering
is available in this Linux workspace. A successful Xcode build is not visual QA.

1. Login: poor hierarchy. Oversized terminal precedes the form; a separate
   three-row progress section creates empty separators and technical SAP text.
   Keep form and compact spinner/status/cancel together; show retry detail only
   when nonempty. Collapse terminal after form and use plain sign-in language.
2. Version loading: a large navigation title and loading card precede a mostly
   empty app section; externalVersionId input is presented before useful results.
   Use inline title, app summary and compact loading/acquisition state. Show
   results as 44-point minimum rows; defer manual ID entry to collapsed Advanced.
3. Main actions: terminal dominates the screen and bare Install/Export links
   conflict with the original full-width PartyUI buttons. Collapse terminal,
   keep Choose version primary and group Install/Export with the actual saved
   app/version label. Preserve favourite styling and the original navigation.
4. Downloads: metadata, actions, details and a prominent destructive action
   compete in one card. Use app/version header, matching Install/Export buttons,
   collapsed Details and delete via labeled menu/swipe with confirmation.

Accessibility risks visible/source-based: ambiguous install target, very small
terminal text, large action stack reducing visible content, spinner-only cues
and destructive action prominence. Add install target labels, accompanying
loading text, 44-point version rows/options target and Dynamic Type fonts.
VoiceOver reading order, text wrapping/contrast, actual spacing, sheet behavior
and touch targets still need device/simulator validation (TESTING.md). No claim
of complete accessibility compliance or verified screenshot fidelity is made.

No broad redesign, web prototype, private UI API or new visual framework.

## Follow-up user screenshots, corrected in 23011

The user reported 23010's main log clipped within a large platter and a meaningless
blue block in the version summary. Source inspection confirmed TerminalPlatter's
250-point outer frame combined with an inner 140-point frame. Remove the inner
frame, preserve original centered container/padding, add full reader with Copy
and optional Follow latest, and remove generic app.fill placeholders.
The screenshot also showed a stale checkmark/acquiring status after rejection;
active Store work now shows a spinner and failures clear that stale state.
The native after-state is still pending device visual verification.
