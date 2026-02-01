# Quest Journal

## Interface

The quest journal is a log of quests and their state for the player
to open and study at any time. The journal should contain information
about currently active quests, but should also hold records of
quests that are no longer active.

Quests are ordered by type. Every category (Main Quest,
Faction Quest, Companion Quest, etc.) can be collapsed individually.

Within a category, the subcategories are collapsible, too, if
applicable. For example, in the Faction Quest category, we might
want to collapse all subcategories and then just open one of them,
such as the one for the Courier Guild.

Quests are be shown with the quest display name, the
current task, and a chronologically sorted list of journal entries
for this particular quest.  For inactive quests, the stage task
is replaced by the quest termination type, such as "Completed",
"Failed" or "Declined".

Inactive quests are hidden by default. A checkbox ("Show inactive quests")
can be toggled to show them or hide them.

## Language

### Quest Names

Quest names are short, unique and to the point. Good examples might
be "Ghost hunt in the attic" or "The seasick captain". It's fine
for quest names to have a bit of artistic freedom, pop culture references
or humor. (e.g. "What could possibly go wrong?" or "Wife of Brian")
Quest names need not convey a clear understanding of what the quest
is about. They just need to be unique, recognizable and maybe memorable.

### Quest Stage Descriptions

Quest stage descriptions indicate what the player is supposed to do next.
They should be plain and clear, since players should remember immediatly
what they were supposed to do when looking at the quest stage description.

Good examples might be "Provide invisibility potion to Mellis" or
"Gather twelve black candles". Sometimes, the current task is vague or
unclear on purpose. In these cases, the stage description should
not provide any spoilers or hints. Stage descriptions can be
lengthy, for example if there are several options on what to do next.

Remember that quest flows can split and follow several avenues in
parallel. In these cases, there might be several stages active at
the same time, and all of their descriptions should be show in the
Quest log - ideally with an indication of how they relate to each
other. (e.g. "Do this OR that", "Do two out of these four", etc.)

### Journal Entries

Journal entries (the normal ones, that aren't quest names or quest
stage tasks) are phrased in past tense. They explain things that
have already happened. (such as "The major wanted you to repair the
fountain. You have accepted the task.")

Journal entries address the player as "you". The entries don't use
plural forms, such as "we", "our" or "us".

When referring to NPCs, journal entries can use either their
name ("Drumin") or their function ("the blacksmith") as approriate.
But keep in mind that the player may not remember who is who in the
entire town, and might therefore need a hint every now and then.
Using only names can be hard when there is more than a hand full of
NPCs.

