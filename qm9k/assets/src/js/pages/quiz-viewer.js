import { closest, mustSingle, mustSingleRef, toggleIf } from "../lib/helpers.js";
import { getJSON, postJSON } from "../lib/post.js";
import { enableVoting, setVoteStatus } from "../components/vote-continue.js";

export function quizMain() {
	let qcont = mustSingle(".-js-quiz-questions .question-container")

	setInterval( updateQuizStatus, 1200 );
	updateQuizStatus();
}

let currentRound = -2;

async function updateQuizStatus() {
	let quizkey = mustSingle("main").dataset["quizkey"];
	let status = await getJSON("/quiz-status/"+quizkey)

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

	toggleIf( mustSingle(".-js-quiz-global-start"), !status.QuizStatus.Started )
	toggleIf( mustSingle(".-js-quiz-questions"), status.QuizStatus.Started && !status.QuizStatus.Grading && !status.QuizStatus.Finished )
	toggleIf( mustSingle(".-js-quiz-grading"), status.QuizStatus.Grading )
	toggleIf( mustSingle(".-js-quiz-global-end"), status.QuizStatus.Finished )


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
	let _ = await postJSON("/set-answer/"+quizkey, {round, question, text});
}
