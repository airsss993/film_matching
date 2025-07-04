// DOM Elements
const navButtons = document.querySelectorAll('.nav-button');
const viewSections = document.querySelectorAll('.view-section');
const logoutButton = document.getElementById('logoutButton');
const likeButton = document.getElementById('likeButton');
const dislikeButton = document.getElementById('dislikeButton');
const matchesGrid = document.getElementById('matchesGrid');
const profileName = document.getElementById('profileName');
const totalSwipes = document.getElementById('totalSwipes');
const totalMatches = document.getElementById('totalMatches');
const preferenceItems = document.getElementById('preferenceItems');

// API Base URL
const API_BASE_URL = 'http://localhost:8080';

// Current film data
let currentFilm = null;

// Navigation
navButtons.forEach(button => {
    button.addEventListener('click', () => {
        const targetView = button.dataset.view;
        
        // Update active states
        navButtons.forEach(btn => btn.classList.remove('active'));
        viewSections.forEach(section => section.classList.remove('active'));
        
        button.classList.add('active');
        document.getElementById(`${targetView}View`).classList.add('active');

        // Load view-specific data
        if (targetView === 'discover') {
            loadNextFilm();
        } else if (targetView === 'matches') {
            loadMatches();
        } else if (targetView === 'profile') {
            loadProfile();
        }
    });
});

// Logout
logoutButton.addEventListener('click', async () => {
    try {
        const response = await fetch(`${API_BASE_URL}/api/auth/logout`, {
            method: 'POST',
            credentials: 'include'
        });

        if (response.ok) {
            window.location.href = 'file:///Users/airsss/Desktop/film_matching/web/auth.html';
        } else {
            throw new Error('Logout failed');
        }
    } catch (error) {
        console.error('Logout error:', error);
    }
});

// Film interactions
likeButton.addEventListener('click', () => handleSwipe(true));
dislikeButton.addEventListener('click', () => handleSwipe(false));

// Load next film
async function loadNextFilm() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/films/next`, {
            credentials: 'include'
        });

        if (!response.ok) throw new Error('Failed to load next film');

        currentFilm = await response.json();
        updateFilmCard(currentFilm);
    } catch (error) {
        console.error('Error loading next film:', error);
    }
}

// Update film card with new data
function updateFilmCard(film) {
    if (!film) return;

    document.getElementById('currentFilmPoster').src = film.poster_url;
    document.getElementById('currentFilmTitle').textContent = film.title;
    document.getElementById('currentFilmYear').textContent = film.year;
    document.getElementById('currentFilmRating').textContent = `★ ${film.rating}`;
    document.getElementById('currentFilmDuration').textContent = `${film.duration} min`;
    document.getElementById('currentFilmDescription').textContent = film.description;

    const genresContainer = document.getElementById('currentFilmGenres');
    genresContainer.innerHTML = '';
    film.genres.forEach(genre => {
        const genreSpan = document.createElement('span');
        genreSpan.className = 'film-genre';
        genreSpan.textContent = genre;
        genresContainer.appendChild(genreSpan);
    });
}

// Handle swipe action
async function handleSwipe(liked) {
    if (!currentFilm) return;

    try {
        const response = await fetch(`${API_BASE_URL}/api/films/swipe`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                film_id: currentFilm.id,
                liked: liked
            }),
            credentials: 'include'
        });

        if (!response.ok) throw new Error('Swipe failed');

        // Check if it's a match
        const result = await response.json();
        if (result.is_match) {
            showMatchAnimation(currentFilm);
        }

        // Load next film
        loadNextFilm();
    } catch (error) {
        console.error('Swipe error:', error);
    }
}

// Show match animation
function showMatchAnimation(film) {
    const matchOverlay = document.createElement('div');
    matchOverlay.className = 'match-overlay';
    matchOverlay.innerHTML = `
        <div class="match-content">
            <h2>It's a Match!</h2>
            <p>You and your friend both liked "${film.title}"</p>
            <img src="${film.poster_url}" alt="${film.title}">
        </div>
    `;

    document.body.appendChild(matchOverlay);

    // Remove overlay after animation
    setTimeout(() => {
        matchOverlay.remove();
    }, 3000);
}

// Load matches
async function loadMatches() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/matches`, {
            credentials: 'include'
        });

        if (!response.ok) throw new Error('Failed to load matches');

        const matches = await response.json();
        updateMatchesGrid(matches);
    } catch (error) {
        console.error('Error loading matches:', error);
    }
}

// Update matches grid
function updateMatchesGrid(matches) {
    matchesGrid.innerHTML = '';

    matches.forEach(match => {
        const matchCard = document.createElement('div');
        matchCard.className = 'match-card';
        matchCard.innerHTML = `
            <div class="match-poster">
                <img src="${match.film.poster_url}" alt="${match.film.title}">
            </div>
            <div class="match-info">
                <h3 class="match-title">${match.film.title}</h3>
                <div class="match-meta">
                    <span>${match.film.year}</span>
                    <span>★ ${match.film.rating}</span>
                </div>
            </div>
        `;
        matchesGrid.appendChild(matchCard);
    });
}

// Load profile data
async function loadProfile() {
    try {
        const response = await fetch(`${API_BASE_URL}/api/users/profile`, {
            credentials: 'include'
        });

        if (!response.ok) throw new Error('Failed to load profile');

        const profile = await response.json();
        updateProfile(profile);
    } catch (error) {
        console.error('Error loading profile:', error);
    }
}

// Update profile information
function updateProfile(profile) {
    profileName.textContent = profile.name;
    totalSwipes.textContent = profile.total_swipes;
    totalMatches.textContent = profile.total_matches;

    // Update preferences
    preferenceItems.innerHTML = '';
    profile.preferences.forEach(preference => {
        const prefSpan = document.createElement('span');
        prefSpan.className = 'preference-item';
        prefSpan.textContent = preference;
        preferenceItems.appendChild(prefSpan);
    });
}

// Add keyboard controls
document.addEventListener('keydown', (event) => {
    if (document.querySelector('#discoverView.active')) {
        if (event.key === 'ArrowLeft') {
            handleSwipe(false); // Dislike
        } else if (event.key === 'ArrowRight') {
            handleSwipe(true); // Like
        }
    }
});

// Initial load
loadNextFilm(); 