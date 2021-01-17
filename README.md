QuizMeister9000 is a web app for running a pub quiz among peers. ("Peers" in the sense that everyone is the quiz master for one round of questions.)

Gameplay
--------
When opening up the app, players can create a profile by answering these questions three.

<p align="center">
<img alt="Creating a profile" src=".readme/qm9k-fig-0.png" width="251" />
</p>

If you have admin rights, the start page features a button that starts a new game.

<p align="center">
<img alt="Admin start page" src=".readme/qm9k-fig-1.png" width="500" />
</p>

This generates an invite code that other players may enter in order to join the quiz.

<p align="center">
<img alt="Joining a quiz using an invite code" src=".readme/qm9k-fig-2.png" width="250" />
</p>

Once everyone's here, players vote to continue the game.

<p align="center">
<img alt="Voting to start the quiz" src=".readme/qm9k-fig-3.png" width="250" />
</p>

In each round, one player (selected randomly) reads their questions, while the others try to answer them. The question description becomes visible on the players' screens while the question-taker is typing.

<p align="center">
<img alt="While typing, other players' screens get updated live" src=".readme/qm9k-fig-4.png" width="250" />
</p>

The contestants' avatars indicate if they're still typing, or have entered text.

<p align="center">
<img alt="Your peers' avatars indicate when they're typing" src=".readme/qm9k-fig-5.png" width="500" />
</p>

After all rounds are over, it's time to check the answers. To prevent favouritism, this is done anonymously.
Answers can be scored as correct, incorrect, or half-right. If multiple people answer the same, it'll need grading only once.

<p align="center">
<img alt="Answers can be correct, incorrect, or 'meh'" src=".readme/qm9k-fig-6.png" width="400" />
</p>

When everyone has graded their papers, the final score is calculated.

<p align="center">
<img alt="The final score" src=".readme/qm9k-fig-7.png" width="400" />
</p>

Building
--------
QuizMeister9000 has some compile-time dependencies:

* Go (≥ 1.15, though some older releases wil probably also work)
* NodeJS (≥ 14 - using a LTS release is highly recommended)

To compile this utility, try running:

    go run build.go

License
-------
This program and its source code are available under the terms of the BSD 3-clause license.
Find out what that means here: https://www.tldrlegal.com/l/bsd3
