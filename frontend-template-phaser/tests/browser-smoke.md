# Browser Smoke Test Checklist

Manual checklist for verifying a Phaser H5 game in the browser before packaging.

## Pre-conditions

- [ ] Game Designer server is running on the expected URL
- [ ] Frontend dev server started with `npm run dev` or preview with `npm run preview`

## Canvas Rendering

- [ ] Page loads without JavaScript errors in console
- [ ] Phaser canvas is visible (not blank/white)
- [ ] Game title text is visible
- [ ] Reel symbols are visible in the 3x3 grid
- [ ] Balance display shows a numeric value
- [ ] Wager display shows a numeric value
- [ ] Spin button is visible and has "SPIN" label

## Controls

- [ ] Spin button responds to click/tap
- [ ] Balance updates after spin
- [ ] Reel symbols change after spin
- [ ] Win message appears after winning spin
- [ ] "No win" message appears after losing spin
- [ ] "Insufficient balance" appears when balance is too low

## Mobile Viewport

- [ ] Game fits within 375x667 viewport (iPhone SE)
- [ ] Game fits within 390x844 viewport (iPhone 14)
- [ ] No horizontal scroll
- [ ] Spin button is easily tappable (at least 44x44px)
- [ ] Text is readable without zooming

## Console Health

- [ ] No JavaScript errors in console during page load
- [ ] No JavaScript errors during spin
- [ ] No 404 errors for assets
- [ ] No CORS errors on API calls

## Error Recovery

- [ ] Game shows "Connection failed" message when server is unreachable
- [ ] Game shows "Spin failed" on API error
- [ ] Game remains usable after error (can retry)
