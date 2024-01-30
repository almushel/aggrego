const lsKeyName = "aggrego-apikey";

function onEnterCallback(e, callback) {
	if (e.key == "Enter") {
		callback();
	}
}

async function login(apikey) {
	const user = await getUser(apikey);
	if (!user) {
		return false;
	}

	document.getElementById("current-user").innerText = "Current user: " + user.name;
	document.getElementById("logged-in").style.display = "block";

	document.getElementById("logged-out").style.display = "none";

	updatePosts(apikey);
	updateFeedList(apikey);

	return true;
}

async function loginWithKeyInput() {
	const keyInput = document.getElementById("apikey-input");

	const success = await login(keyInput.value)
	if (success) {
		localStorage.setItem(lsKeyName, keyInput.value)
	}
}

function logout() {
	localStorage.removeItem(lsKeyName)

	document.getElementById("logged-in").style.display = "none";
	document.getElementById("posts-container").textContent = "";
	document.getElementById("logged-out").style.display = "block";
}

function loadPage(e) {
	e.preventDefault();
	const page = e.target.href;
	window.history.pushState({}, "", page);
	lsUpdatePosts();
}

function newPostElement(post) {
	const result = document.createElement("div")
	const title = document.createElement("h3");
	const url = document.createElement("a");
	url.href=post.url;
	url.innerText = post.title;

	const description = document.createElement("p");
	description.innerHTML = post.description;

	title.append(url);
	result.append(title, description);
	return result;
}

async function updatePosts(apikey) {
	const posts = await getPosts(apikey);
	const container = document.getElementById("posts-container");
	const children = [];
	if (posts.list) {
		for (let post of posts.list) {
			children.push(newPostElement(post))
		}
	} else {
		const notFound = document.createElement("p")
		notFound.textContent = `Page ${posts.page} not found`;
		children.push(notFound);
	}

	const pageCount = Math.ceil(posts.totalCount/pageSize);

	const pages = document.createElement("div");
	pages.style.display = "flex";
	pages.style.justifyContent = "center";
	for (let pg = 1; pg <= pageCount; pg++) {
		const pgLink = document.createElement("a")
		const separator = pg != pageCount ? "," : ""
		if (pg != posts.page) {
			pgLink.addEventListener("click", loadPage);
			pgLink.href = `${window.location.pathname}?pg=${pg}`;
			pgLink.textContent = `  ${pg} ${separator}`;
		} else {
			pgLink.innerHTML = `<strong>  ${pg} ${separator}</strong>`;
		}
		
		pages.append(pgLink);
	}

	children.push(pages);
	container.replaceChildren(...children);
}

async function updateFeedList(apikey) {
	const feeds = await getFeeds();
	const follows = await getFeedFollows(apikey);

	const followIDs = []
	for (let f of follows) {
		followIDs.push(f.feed_id)	
	}

	const feedList = document.getElementById("user-feeds");

	let count = 1;
	let elements = []
	for (let feed of feeds) {
		const url = new URL(feed.url);
		const elementID = `feed-list-item-${count}`;

		const e = document.createElement("li");
		e.innerHTML = 
			`<input type="checkbox" id="${elementID}" ${followIDs.includes(feed.id) ? "checked" : ""}>` +
			`<label for="${elementID}">${url.hostname}</label>`

		const checkbox = e.children[elementID];
		checkbox.addEventListener("input", async () => {
			const success = await toggleFeedFollow(feed.id, checkbox.checked)
			if (success) {
				lsUpdatePosts()
			} else {
				checkbox.checked = !checkbox.checked;
			}
		})

		elements.push(e)
		count++;
	}

	feedList.replaceChildren(...elements)
}

async function submitNewFeed() {
	const apikey = localStorage.getItem(lsKeyName)
	if (apikey == null) {
		return false;
	}
	const nameField = document.getElementById("new-feed-name");
	const urlField = document.getElementById("new-feed-url");

	if (nameField.value && urlField.value) {
		const success = await postFeed(nameField.value, urlField.value, apikey)
		if (success) {
			nameField.value = "";
			urlField.value = "";

			updateFeedList(apikey);
			updatePosts(apikey);
		}
	}
}

async function lsLogin() {
	const apikey = localStorage.getItem(lsKeyName)
	if (apikey == null) {
		return false;
	}
	login(apikey)
}

async function lsUpdatePosts() {
	const apikey = localStorage.getItem(lsKeyName)
	if (apikey == null) {
		return false;
	}
	updatePosts(apikey)
}