package main

import (
	"math/rand"
	"time"
)

var(
	QUOTES = [...]string{
		"Have the courage to follow your heart and your intuition. They somehow already know what you truly want to become. - Steve Jobs",
		"Perfection is attained, not when there is nothing more to add, but when there is nothing more to take away. - Antoine de Saint Exupéry, French writer",
		"After love, the sense of productivity is the most fulfilling emotion.",
		"One day you’ll understand that it’s harder to be kind than clever. - Jeff Bezo's Grandpa",
		"(F)ollow (O)ne (C)ourse (U)ntil (S)uccess",
		"The master has failed more times than the beginner has tried.",
		"I fear not the man who has practiced 10,000 kicks once, but I fear the man who has practiced one kick 10,000 times.",
		"An amateur practices until he can play it correctly, a professional practices until he can't play it incorrectly.",
		"No one cares who you are until they care about the reason why you need to exist.",
		"Gratitude is the willingness to recognize the unearned increments of value in one's experience.",
		"True literature can exist only where it is created, not by diligent and trustworthy functionaries, but by madmen, hermits, heretics, dreamers, rebels, and skeptics. Yevgeny Zamyatin, A Soviet Heretic: Essays by Yevgeny Zamyatin",
		"Technology is all about deviating from the natural. Figuring out how to manipulate and distort the natural in a way that works for us — Aubrey de Gray",
		"Maturity is the ability to reject good alternatives in order to pursue even better ones. - Ray Dalio",
		"People only take into account two things. the peak of how good it was, and how it ended.",
		"The easiest way out is through. Don't fight if you can't win. Be stoic and it'll be over soon enough.",
		"The saddest people smile the brightest. All because they do not wish to see anyone suffer the way they do. - Unknown",
		"The way to get started is to quit talking and begin doing. - Walt Disney",
		"Psychologists have found that we are each more interested in knowing that the other person is trying to empathize with us … than we are in believing that they have actually accomplished that goal. Good listening … is profoundly communicative. And struggling to understand communicates the most positive message of all. —Difficult Conversations, Douglas Stone, Bruce Patton, and Sheila Heen",
		"People will remember not what you said but how you made them feel.",
		"Good leaders are trailblazers, making a path for others to follow. Great leaders, however, inspire their people to reach higher, dream bigger, and achieve greater. Perhaps the most important leadership skill you can develop is the ability to provide inspiration to your team.",
		"A leader is best when people barely know he exists, when his work is done, his aim fulfilled, they will say: we did it ourselves. — Lao Tzu",
		"Man is made by his belief. As he believes, so he is. - The Bhagavad Gita",
		"You know what the difference between a vision and a hallucination is? They call it a vision when other people can see it.",
		"Always enter the conversation that’s already being had in the person’s heart and mind.",
		"Don’t practice what you don’t want to become. Be careful with what you say or do. It will become chemically metastasized. If you did, then you can only inhibit them but they’re permanently hardwired. (paraphrased) - Jordan Peterson",
		"The mind creates the abyss, the heart crosses it. - Nisargadatta Maharaj",
		"Do unto others 20% better than you would expect them to do unto you, to correct for subjective error. - Linus Pauling",
		"Marketing in its essence, is the articulation of value. - Ryan Deiss",
		"Pain is inevitable, suffering is optional. - Haruki Murakami",
		"Good habits come from thinking repeatedly in a principled way. - Ray Dalio",
		"Learning is the product of a continuous real-time feedback loop in which we make decisions, see their outcomes, and improve our understanding of reality as a result. - Ray Dalio",
		"I am not a product of my circumstances. I am a product of my decisions. - Stephen Covey",
		"If someone loves you, you will know. If they don’t, you will be confused.",
		"There will be demands upon your ability, upon your endurance, upon your disposition, upon your patience...just as fire tempers iron into fine steel so does adversity temper one’s character into firmness, tolerance and determination. — Senator Margaret Chase Smith, Lieutenant Colonel",
		"It's better to do something and be wrong than do nothing and be wrong.",
		"Art is a reflection of humanity, and humanity’s greatest virtue is its ability to overcome adversity. - Christopher Zara",
		"Necessity is the mother of invention.",
		"I’ve always measured success not only by how many people love it but by how many people hate it too, because if you do something that people love, the hate is not really worth much. If you don’t do anything people love or hate, you have failed. - Marilyn Manson",
		"Any man who must say ‘I am the king’ is no true king. - Tywin Lannister",
		"The perfect distance for doing nothing is when you have the constant chance to do something.",
		"Projection is a subtle cognitive flaw. Want to understand someone? Ask them to guess other people's motives, and you’ll see their own.",
		"If I skip practice for one day, I notice. If I skip practice for two days, my wife notices. If I skip practice for three days, the world notices. - Vladimir Horowitz",
		"Do what you do best, delegate the rest. - John Forsezee, a mobster",
		"No man ever steps in the same river twice, for it is not the same river and he is not the same man. - Heraclitus",
		"A ship is always safe at the shore - but that is NOT what it is built for. - Albert Einstein",
		"Routines are great, but they can also be a rut - Jesse Itzler",
		"Spontaneous Trait Transference - If you speak ill of someone, people can't help but associate those negative traits with you.",
		"Problem-solving is as much a mindset as it is a collection of discrete skills and methods. No single person, regardless of how talented, will be able to solve every problem. A good problem solver learns how to leverage their unique experiences and skills against hard problems and how to work with others to achieve the best possible result.",
		"We suffer more often in imagination than in reality. - Seneca",
		"Whether you think you can or you can’t - either way, you’re right. - Henry Ford",
		"What you reveal, you heal. - Jay Z",
		"You don't learn to walk by following rules. You learn by doing, and by falling over. - Richard Branson",
		"There is no end to the adventures that we can have if only we seek them with our eyes open. - Jawaharlal Nehru",
		"It is during our darkest moments that we must focus to see the light. - Aristotle",
		"Discipline and concentration are a matter of being interested. Find pleasure in what you are doing. - Tom Kite",
		"The way to get started is to quit talking and begin doing. - Walt Disney",
		"One way to remember who you are is to remember who your heroes are. - Walter Isaacson",
		"The art of being wise is the art of knowing what to overlook. - William James",
		"You can’t use up creativity. The more you use, the more you have. - Maya Angelou",
		"Motivation is what gets you started. Habit is what keeps you going. - Jim Ryun",
		"Everything you’ve ever wanted is on the other side of fear. - George Addair",
	}
)

func GetQuote() string {
	rand.Seed(time.Now().Unix())
	return QUOTES[rand.Intn(len(QUOTES))]
}