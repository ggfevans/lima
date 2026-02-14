Output demo.gif

Set Shell "zsh"
Set Theme "Dracula"

# Larger canvas to avoid wrapping in README embeds.
Set FontSize 28
Set Width 1400
Set Height 820
Set Padding 20
Set WindowBar "Rings"

# Slight overall speed-up while keeping sleeps reasonably robust.
Set Framerate 30
Set PlaybackSpeed 1.2
Set TypingSpeed 25ms

Sleep 800ms

# Run the demo.
Type "../bin/endorse --demo" Enter
Sleep 2.5s

# Step through the prompts.
Enter
Sleep 2s
Type "uhh" Enter
Sleep 2s
Type "you should leave" Enter
Sleep 2.5s

# Navigation (recorded as escape sequences originally).
Up 2
Left
Down
Enter
Sleep 2s

Type "what" Enter
Sleep 2.5s
Type "r u crazy" Enter
Sleep 2.5s

Up
Left
Down
Sleep 800ms
Enter
Sleep 600ms

Right 2
Sleep 800ms
Enter
Sleep 600ms

Left
Down 2
Sleep 700ms
Enter
Sleep 1.8s

Type "im in" Enter
Sleep 2.5s
Type "amazing" Enter
Sleep 1.2s

Left 2
Enter
Down
Left 2
Down
Sleep 1.2s
Enter
Sleep 1.8s

Type "I am telling HR" Enter
Sleep 900ms

Right 2
Sleep 800ms
Type "y"
Sleep 300ms
Type "no"
Right
Sleep 500ms
Type "no!!!!" Enter
Sleep 1.5s

Ctrl+C
