import { all, closest, mustSingle, mustSingleRef, onClick, toggleIf } from "../lib/helpers.js";
import { getJSON, postJSON } from "../lib/post.js";
import { enableVoting, setVoteStatus } from "../components/vote-continue.js";

export function quizMain() {
	let qcont = mustSingle(".-js-quiz-questions .question-container")

	setInterval( updateQuizStatus, 1200 );
	updateQuizStatus();

	// setInterval( updateGrading, 600 );
	onClick(".-js-click-grade", applyGrade);
}

let currentRound = -2;
let wasGrading = false;
let amGrading = false;

let wasFinished = false;

async function updateQuizStatus() {
	let quizkey = mustSingle("main").dataset["quizkey"];
	let status = await getJSON("quiz-status/"+quizkey)

	enableVoting(status.VotingEnabled);
	setVoteStatus(status.MyVote);

	if ( status.QuizStatus.Started && !status.QuizStatus.Grading && !status.QuizStatus.Finished ) {
		let qcont = mustSingle(".-js-quiz-questions .question-container")

		if ( status.CurrentRound.RoundNo != currentRound ) {
			// The round has changed. Update the form
			currentRound = status.CurrentRound.RoundNo;

			qcont.innerHTML = "";

			let template = mustSingle(".-js-template-answer");
			let tClass = "answer";
			if ( status.CurrentRound.ThisIsMe ) {
				template = mustSingle(".-js-template-question");
				tClass = "question";
			}
			status.CurrentRound.Questions.forEach(q => {
				let ndq = document.createElement("DIV");
				ndq.classList.add(tClass);
				ndq.innerHTML = template.innerHTML;

				if ( status.CurrentRound.ThisIsMe ) {
					mustSingleRef(ndq, "textarea").value = q.Question;
				} else {
					mustSingleRef(ndq, "-question").innerText = q.Question;
					mustSingleRef(ndq, "textarea").value = q.MyAnswer;
				}

				let txt = mustSingleRef(ndq, "textarea");
				txt.addEventListener("focus",questionFocus);
				txt.addEventListener("blur",questionBlur);

				qcont.appendChild(ndq);
				mustSingleRef(ndq, ".-number").innerText = qcont.children.length;
			});

			let title = mustSingle(".-js-quiz-questions h3")
			if ( status.CurrentRound.ThisIsMe ) {
				title.innerHTML = "This is you.";
			} else {
				title.innerHTML = "Please direct your attention to: ";
				let name = document.createElement("STRONG");
				name.innerText = status.CurrentRound.QuizMaster.Nick;
				title.appendChild(name);
			}

			window.setTimeout(() => { mustSingleRef(qcont, "textarea").focus(); }, 20);
		}

		let nQuestions = qcont.children.length;

		if ( status.CurrentRound.ThisIsMe ) {
			qcont.querySelectorAll(".question").forEach((ndq, i) => {
				const q = status.CurrentRound.Questions[i];
				if ( !q ) {
					return;
				}

				// Check textareas for changes, and save if they're different
				let txt = mustSingleRef(ndq,"textarea");
				if ( txt.value != q.Question ) {
					setTimeout(() => { setAnswer( currentRound, i, txt.value ) }, 50*(nQuestions-i));
				}
			})
		} else {
			qcont.querySelectorAll(".answer").forEach((ndq, i) => {
				const q = status.CurrentRound.Questions[i];
				if ( !q ) {
					return;
				}

				mustSingleRef(ndq, ".-question").innerText = q.Question;

				// Check textareas for changes, and save if they're different
				let txt = mustSingleRef(ndq,"textarea");
				if ( txt.value != q.MyAnswer ) {
					setTimeout(() => { setAnswer( currentRound, i, txt.value ) }, 50*(nQuestions-i));
				}
			})
		}
	}
	if ( status.QuizStatus.Grading ) {
		if ( !wasGrading ) {
			wasGrading = true;

			let gcont = mustSingle(".grading-container");
			gcont.innerHTML = "";

			let grading = await getJSON("grade-answers/"+quizkey);
			grading.Questions.forEach((q,i) => {
				let ndq = document.createElement("DIV");
				ndq.classList.add("question");
				ndq.innerHTML = mustSingle(".-js-quiz-grading .-js-template-question").innerHTML;
				ndq.dataset["question"] = i;

				mustSingleRef(ndq, ".-js-question").innerText = (i+1) + ": " + q.Question;

				q.Answers.forEach(ans => {
					let nda = document.createElement("DIV");
					nda.classList.add("answer");
					nda.innerHTML = mustSingle(".-js-quiz-grading .-js-template-answer").innerHTML;
					nda.dataset["answer"] = ans.Answer;
					mustSingleRef(nda, ".-js-answer").innerText = ans.Answer;
					ndq.appendChild(nda);
				});

				gcont.appendChild(ndq);
			});

			amGrading = true;
			readGrading(grading);
		}
	}

	if ( status.QuizStatus.Finished ) {
		if ( !wasFinished ) {
			wasFinished = true;

			let lb = mustSingle(".leaderboard-container");
			lb.innerHTML = "";

			let leaderboard = await getJSON("leaderboard/"+quizkey);

			let lastScore = 0, lastPos = 0;

			leaderboard.Peers.forEach((p,i) => {
				let ndc = document.createElement("DIV");
				ndc.classList.add("contestant");
				ndc.innerHTML = mustSingle(".-js-quiz-global-end .-js-template-contestant").innerHTML;

				let pos = i+1;
				if ( p.Score == lastScore ) {
					pos = lastPos;
				}
				lastScore =  p.Score;
				lastPos = pos;

				mustSingleRef(ndc, ".-nick").innerText = `#${pos}: ${p.Nick}`;
				mustSingleRef(ndc, ".-quest").innerText = `Quest: ${p.Quest}`;
				mustSingleRef(ndc, ".-avatar svg use").setAttribute("fill", `#${p.Colour}`);

				let score = Math.floor(p.Score).toString();
				if ( p.Score % 2 == 1 ) {
					score += "½";
				}
				mustSingleRef(ndc, ".-score").innerText = `Score: ${score}`;

				lb.appendChild(ndc);
			})
		}
	}

	toggleIf( mustSingle(".-js-quiz-global-start"), !status.QuizStatus.Started )
	toggleIf( mustSingle(".-js-quiz-questions"), status.QuizStatus.Started && !status.QuizStatus.Grading && !status.QuizStatus.Finished )
	toggleIf( mustSingle(".-js-quiz-grading"), status.QuizStatus.Grading )
	toggleIf( mustSingle(".-js-quiz-global-end"), status.QuizStatus.Finished )
	toggleIf( mustSingle(".-js-global-peer-status"), !status.QuizStatus.Finished )
}

function questionFocus(e) {
	let nd = closest( e.target, ".question,.answer" );
	if ( nd ) {
		nd.classList.toggle("-focus", true);
	}
}
function questionBlur(e) {
	let nd = closest( e.target, ".question,.answer" );
	if ( nd ) {
		nd.classList.toggle("-focus", false);
	}
}

async function setAnswer(round, question, text) {
	let quizkey = mustSingle("main").dataset["quizkey"];
	let _ = await postJSON("set-answer/"+quizkey, {round, question, text});
}


async function updateGrading() {
	if ( !amGrading ) {
		return;
	}

	let quizkey = mustSingle("main").dataset["quizkey"];
	let grading = await getJSON("grade-answers/"+quizkey);
	readGrading(grading);
}

async function applyGrade(e) {
	let score = e.target.dataset["score"];
	let nda = closest(e.target, ".answer");
	let answer = nda.dataset["answer"];
	let ndq = closest(nda, ".question");
	let question = ndq.dataset["question"];

	let quizkey = mustSingle("main").dataset["quizkey"];
	let grading = await postJSON("grade-answers/"+quizkey, {question, answer, score});
	readGrading(grading);
}

function readGrading(grading) {
	let ndqs = all(".grading-container .question");
	ndqs.forEach((ndq,i) => {
		ndq.querySelectorAll(".answer").forEach(nda => {
			let answer = nda.dataset["answer"];
			grading.Questions[i].Answers.forEach(ans => {
				if ( ans.Answer != answer ) {
					return;
				}

				nda.querySelectorAll(".-grade .tfbutton").forEach(b => {
					let selected = ans.Scored && (b.dataset["score"] == ans.Score);
					b.classList.toggle("-white", !selected);
				});
			});
		});
	});
}
