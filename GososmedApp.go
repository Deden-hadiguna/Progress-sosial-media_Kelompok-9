package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	MAXUSER    = 100
	MAXPOST    = 1000
	MAXCOMMENT = 5000
	MAXFRIEND  = 50
)

// Struktur Data
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Nama      string `json:"nama"`
	Bio       string `json:"bio"`
	JoinDate  string `json:"join_date"`
	PostCount int    `json:"post_count"`
	Friends   []int  `json:"friends"`
}

type Post struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	Timestamp string    `json:"timestamp"`
	Likes     []int     `json:"likes"`
	Comments  []Comment `json:"comments"`
}

type Comment struct {
	ID        int    `json:"id"`
	PostID    int    `json:"post_id"`
	UserID    int    `json:"user_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
}

// Variabel Global
var users []User
var posts []Post
var currentUser int

const (
	userFile = "users.json"
	postFile = "posts.json"
)

// Fungsi Utility
func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func showHeader(title string) {
	fmt.Println("\n=================================")
	fmt.Println("🌟 GOsosmedAPP v1.0")
	fmt.Println("=================================")
	fmt.Println("Kelompok 9")
	fmt.Println("=================================")
	fmt.Printf("📍 %s\n", title)
	fmt.Println("=================================")
}

// File Operations
func saveToFile(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, jsonData, 0644)
}

func loadFromFile(filename string, data interface{}) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}
	fileData, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileData, data)
}

func saveUsersToFile() {
	if err := saveToFile(users, userFile); err != nil {
		fmt.Println("❌ Gagal menyimpan data pengguna:", err)
	}
}

func savePostsToFile() {
	if err := saveToFile(posts, postFile); err != nil {
		fmt.Println("❌ Gagal menyimpan data post:", err)
	}
}

func loadUsersFromFile() {
	if err := loadFromFile(userFile, &users); err != nil {
		fmt.Println("❌ Gagal memuat data pengguna:", err)
	}
}

func loadPostsFromFile() {
	if err := loadFromFile(postFile, &posts); err != nil {
		fmt.Println("❌ Gagal memuat data post:", err)
	}
}

func getUserById(id int) *User {
	for i := range users {
		if users[i].ID == id {
			return &users[i]
		}
	}
	return nil
}

// User Management
func showProfile() {
	showHeader("My Profile")
	user := getUserById(currentUser)
	if user == nil {
		fmt.Println("❌ Gagal memuat profil pengguna.")
		return
	}

	fmt.Printf("👤 Nama: %s\n", user.Nama)
	fmt.Printf("🆔 Username: %s\n", user.Username)
	fmt.Printf("📌 Bio: %s\n", user.Bio)
	fmt.Printf("📅 Join Date: %s\n", user.JoinDate)
	fmt.Printf("👥 Friends: %d\n", len(user.Friends))
	fmt.Printf("📝 Posts: %d\n", user.PostCount)
	fmt.Println("==============================")
	fmt.Println("Daftar Teman:")
	for _, friendID := range user.Friends {
		friend := getUserById(friendID)
		if friend != nil {
			fmt.Printf("- %s (%s)\n", friend.Nama, friend.Username)
		}
	}
}

func login() bool {
	showHeader("Login")

	username := getInput("👤 Username: ")
	password := getInput("🔑 Password: ")

	for _, user := range users {
		if user.Username == username && user.Password == password {
			currentUser = user.ID
			fmt.Printf("\n✨ Welcome back, %s!\n", user.Nama)
			return true
		}
	}

	fmt.Println("\n❌ Username atau password salah!")
	fmt.Println("\n1. Coba Lagi")
	fmt.Println("2. Kembali ke Menu Utama")

	choice := getInput("\n🎯 Pilihan: ")

	if choice == "1" {
		return login()
	}

	return false
}

func registerUser() {
	showHeader("Register New User")

	username := getInput("📝 Username: ")

	// Cek username unik
	for _, user := range users {
		if user.Username == username {
			fmt.Println("❌ Username sudah dipakai!")
			return
		}
	}

	password := getInput("🔑 Password: ")
	nama := getInput("👤 Nama Lengkap: ")
	bio := getInput("📌 Bio: ")

	newUser := User{
		ID:       len(users) + 1,
		Username: username,
		Password: password,
		Nama:     nama,
		Bio:      bio,
		JoinDate: time.Now().Format("2006-01-02"),
		Friends:  []int{},
	}

	users = append(users, newUser)
	saveUsersToFile()

	fmt.Println("\n✅ Registrasi berhasil! Silahkan login.")
}

func editProfile() {
	showHeader("Edit Profile")
	user := getUserById(currentUser)
	if user == nil {
		fmt.Println("❌ Gagal mengakses profil")
		return
	}

	fmt.Println("1. Edit Nama")
	fmt.Println("2. Edit Bio")
	fmt.Println("3. Edit Password")
	fmt.Println("4. Kembali")

	choice := getInput("\n🎯 Pilihan: ")

	switch choice {
	case "1":
		user.Nama = getInput("Nama Baru: ")
	case "2":
		user.Bio = getInput("Bio Baru: ")
	case "3":
		user.Password = getInput("Password Baru: ")
	case "4":
		return
	default:
		fmt.Println("❌ Pilihan tidak valid")
		return
	}

	saveUsersToFile()
	fmt.Println("✅ Profil berhasil diupdate!")
}

// Timeline & Posts
func showTimeline() {
	showHeader("Timeline")

	if len(posts) == 0 {
		fmt.Println("📭 Belum ada postingan...")
		return
	}

	for i := len(posts) - 1; i >= 0; i-- {
		post := posts[i]
		user := getUserById(post.UserID)
		if user != nil {
			fmt.Println("\n--------------------------------")
			fmt.Printf("👤 %s\n", user.Nama)
			fmt.Printf("🕒 %s\n", post.Timestamp)
			fmt.Printf("📝 %s\n", post.Content)
			fmt.Printf("👍 %d Likes | 💬 %d Comments\n", len(post.Likes), len(post.Comments))

			// Show latest comments
			fmt.Println("--- Komentar Terbaru ---")
			showComments(post.Comments)

			// Show interaction menu
			fmt.Println("\n1. Like/Unlike")
			fmt.Println("2. Tambah Komentar")
			fmt.Println("3. Lanjut")

			choice := getInput("Pilihan: ")
			switch choice {
			case "1":
				toggleLike(i)
			case "2":
				addComment(i)
			}
		}
	}
}

func toggleLike(postIndex int) {
	liked := false
	for i, userID := range posts[postIndex].Likes {
		if userID == currentUser {
			// Unlike
			posts[postIndex].Likes = append(posts[postIndex].Likes[:i], posts[postIndex].Likes[i+1:]...)
			liked = true
			break
		}
	}

	if !liked {
		// Like
		posts[postIndex].Likes = append(posts[postIndex].Likes, currentUser)
	}

	savePostsToFile()
	fmt.Println("✅ Like berhasil diupdate!")
}

func addComment(postIndex int) {
	content := getInput("💬 Komentar: ")

	comment := Comment{
		ID:        len(posts[postIndex].Comments) + 1,
		PostID:    posts[postIndex].ID,
		UserID:    currentUser,
		Content:   content,
		Timestamp: time.Now().Format("2006-01-02 15:04"),
	}

	posts[postIndex].Comments = append(posts[postIndex].Comments, comment)
	savePostsToFile()
	fmt.Println("✅ Komentar berhasil ditambahkan!")
}

func showComments(comments []Comment) {
	for i := len(comments) - 1; i >= max(0, len(comments)-3); i-- {
		comment := comments[i]
		user := getUserById(comment.UserID)
		if user != nil {
			fmt.Printf("💬 %s: %s\n", user.Nama, comment.Content)
		}
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func createPost() {
	showHeader("Create New Post")

	if len(posts) >= MAXPOST {
		fmt.Println("❌ Storage post penuh!")
		return
	}

	content := getInput("📝 Apa yang anda pikirkan? ")

	newPost := Post{
		ID:        len(posts) + 1,
		UserID:    currentUser,
		Content:   content,
		Timestamp: time.Now().Format("2006-01-02 15:04"),
		Likes:     []int{},
		Comments:  []Comment{},
	}

	posts = append(posts, newPost)

	user := getUserById(currentUser)
	if user != nil {
		user.PostCount++
		saveUsersToFile()
	}

	savePostsToFile()
	fmt.Println("\n✅ Post berhasil dibuat!")
}

// Friend Management
func showFriends() {
	showHeader("Daftar Teman")

	user := getUserById(currentUser)
	if user == nil || len(user.Friends) == 0 {
		fmt.Println("😢 Anda belum memiliki teman...")
		return
	}

	fmt.Println("Urutkan berdasarkan:")
	fmt.Println("1. Nama (A-Z)")
	fmt.Println("2. Tanggal Bergabung")
	fmt.Println("3. Jumlah Post")

	choice := getInput("\n🎯 Pilihan: ")

	var friends []User
	for _, friendID := range user.Friends {
		if friend := getUserById(friendID); friend != nil {
			friends = append(friends, *friend)
		}
	}

	switch choice {
	case "1":
		sort.Slice(friends, func(i, j int) bool {
			return friends[i].Nama < friends[j].Nama
		})
	case "2":
		sort.Slice(friends, func(i, j int) bool {
			return friends[i].JoinDate < friends[j].JoinDate
		})
	case "3":
		sort.Slice(friends, func(i, j int) bool {
			return friends[i].PostCount > friends[j].PostCount
		})
	}

	fmt.Printf("\nTotal Teman: %d\n\n", len(friends))
	for i, friend := range friends {
		fmt.Printf("%d. %s (%s)\n", i+1, friend.Nama, friend.Username)
		fmt.Printf("   📌 %s\n", friend.Bio)
		fmt.Printf("   🤝 Bergabung: %s\n", friend.JoinDate)
		fmt.Printf("   📝 Jumlah Post: %d\n\n", friend.PostCount)
	}
}

func manageFriends() {
	showHeader("Manajemen Teman")

	fmt.Println("1. Tambah Teman")
	fmt.Println("2. Hapus Teman")
	fmt.Println("3. Kembali")

	choice := getInput("\n🎯 Pilihan: ")

	switch choice {
	case "1":
		addFriend()
	case "2":
		removeFriend()
	}
}

func addFriend() {
	showHeader("Tambah Teman")

	user := getUserById(currentUser)
	if user == nil {
		fmt.Println("❌ Gagal mengakses data pengguna.")
		return
	}

	username := getInput("👤 Username teman: ")

	foundUser := -1
	for _, u := range users {
		if u.Username == username {
			foundUser = u.ID
			break
		}
	}

	if foundUser == -1 {
		fmt.Println("❌ User tidak ditemukan!")
		return
	}

	if foundUser == currentUser {
		fmt.Println("❌ Anda tidak dapat menambahkan diri sendiri sebagai teman!")
		return
	}

	for _, friendID := range user.Friends {
		if friendID == foundUser {
			fmt.Println("❌ Sudah berteman!")
			return
		}
	}

	user.Friends = append(user.Friends, foundUser)
	foundUserRef := getUserById(foundUser)
	if foundUserRef != nil {
		foundUserRef.Friends = append(foundUserRef.Friends, currentUser)
	}
	saveUsersToFile()
	fmt.Printf("\n✅ Berhasil berteman dengan %s!\n", getUserById(foundUser).Nama)
}

func removeFriend() {
	showHeader("Hapus Teman")

	user := getUserById(currentUser)
	if user == nil || len(user.Friends) == 0 {
		fmt.Println("😢 Anda belum memiliki teman...")
		return
	}

	fmt.Println("Daftar Teman:")
	for i, friendID := range user.Friends {
		if friend := getUserById(friendID); friend != nil {
			fmt.Printf("%d. %s (%s)\n", i+1, friend.Nama, friend.Username)
		}
	}

	choice := getInput("\nPilih nomor teman yang akan dihapus (0 untuk batal): ")
	index := atoi(choice) - 1

	if index < 0 || index >= len(user.Friends) {
		return
	}

	friendID := user.Friends[index]
	friend := getUserById(friendID)

	// Hapus dari daftar teman user
	user.Friends = append(user.Friends[:index], user.Friends[index+1:]...)

	// Hapus dari daftar teman friend
	if friend != nil {
		for i, id := range friend.Friends {
			if id == currentUser {
				friend.Friends = append(friend.Friends[:i], friend.Friends[i+1:]...)
				break
			}
		}
	}

	saveUsersToFile()
	fmt.Printf("\n✅ Berhasil menghapus %s dari daftar teman!\n", friend.Nama)
}

func atoi(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

// Menu & Main
func showMainMenu() {
	for {
		showHeader("Main Menu")
		fmt.Println("1. 📱 Timeline")
		fmt.Println("2. ✍️  Create Post")
		fmt.Println("3. 👥 Friends List")
		fmt.Println("4. 🤝 Manage Friends")
		fmt.Println("5. 👤 My Profile")
		fmt.Println("6. ✏️  Edit Profile")
		fmt.Println("7. 🚪 Logout")

		choice := getInput("\n🎯 Pilihan: ")

		switch choice {
		case "1":
			showTimeline()
		case "2":
			createPost()
		case "3":
			showFriends()
		case "4":
			manageFriends()
		case "5":
			showProfile()
		case "6":
			editProfile()
		case "7":
			currentUser = 0
			return
		}
	}
}

func main() {
	loadUsersFromFile()
	loadPostsFromFile()

	for {
		showHeader("Welcome")
		fmt.Println("1. 🔑 Login")
		fmt.Println("2. ✨ Register")
		fmt.Println("3. ❌ Exit")

		choice := getInput("\n🎯 Pilihan: ")

		switch choice {
		case "1":  
			if login() {
				showMainMenu()
			}
		case "2":
			registerUser()
		case "3":
			fmt.Println("👋 Terima kasih telah menggunakan aplikasi!")
			return
		default:
			fmt.Println("❌ Pilihan tidak valid!")
		}
	}
}
