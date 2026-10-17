# 🖼️ The gallery

`cmd/gallery` has a screen for each part of the catalogue, one file each (`-screen 0` to `4` opens on
one, for captures):

<table>
<tr><th>Actions</th><th>Inputs</th><th>Lists</th><th>Feedback</th><th>Tabs</th></tr>
<tr>
<td><picture><source media="(prefers-color-scheme: dark)" srcset="../images/actions-dark.png"><img alt="Buttons, a split button, a button group, assist chips, a floating toolbar" src="../images/actions-light.png" width="150"></picture></td>
<td><picture><source media="(prefers-color-scheme: dark)" srcset="../images/inputs-dark.png"><img alt="A search bar, text fields, a dropdown and a password field" src="../images/inputs-light.png" width="150"></picture></td>
<td><picture><source media="(prefers-color-scheme: dark)" srcset="../images/lists-dark.png"><img alt="A banner over a list of rows with photos and switches" src="../images/lists-light.png" width="150"></picture></td>
<td><picture><source media="(prefers-color-scheme: dark)" srcset="../images/feedback-dark.png"><img alt="Loading indicators, wavy progress and skeletons" src="../images/feedback-light.png" width="150"></picture></td>
<td><picture><source media="(prefers-color-scheme: dark)" srcset="../images/tabs-dark.png"><img alt="Tabs over swiping pages and carousels of photos of a coast" src="../images/tabs-light.png" width="150"></picture></td>
</tr>
</table>

It was tried on a Pixel 6 Pro (Android 16, the wallpaper's colours, dark), the Android 16 emulator
(light, turned all the way round) and a OnePlus 7 Pro (Android 10, Material's baseline): every
screen, the gestures, the keyboard, a notification, going to the background and back. Its test
(`cmd/gallery`) touches every destination from every other with a finger; the material package's
tests drive each piece with a fake pointer and fake fingers, with no display.

---

[← All the documentation](../1-overview.md)
