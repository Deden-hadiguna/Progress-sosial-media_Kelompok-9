package main

import "fmt"

type User struct {
	ID       int
	Username string
	Password string
	Nama     string
}

// Login pakai ini sementara untuk menggunakan aplikasi
var users = []User{
	{ID: 1, Username: "user1", Password: "pass1", Nama: "User One"},
	{ID: 2, Username: "user2", Password: "pass2", Nama: "User Two"},
}

var currentUser *User

func showHeader(title string) {
	fmt.Println("\n=================================")
	fmt.Println("🌟 GOsosmedAPP v1.0")
	fmt.Println("=================================")
	fmt.Println("Kelompok 9")
	fmt.Println("=================================")
	fmt.Printf("📍 %s\n", title)
	fmt.Println("=================================")
}

func login() bool {
	var username, password string
	fmt.Print("👤 Username: ")
	fmt.Scan(&username)
	fmt.Print("🔑 Password: ")
	fmt.Scan(&password)

	for i := range users {
		if users[i].Username == username && users[i].Password == password {
			currentUser = &users[i]
			fmt.Printf("\n✨ Welcome back, %s!\n", currentUser.Nama)
			return true
		}
	}

	fmt.Println("❌ Username atau password salah!")
	return false
}

func showMainMenu() {
	for {
		showHeader("Main Menu")
		fmt.Println("1. 📱 Beranda")
		fmt.Println("2. ✍️  Apa yang anda pikirkan? Posting sesuatu disini")
		fmt.Println("3. 👥 Daftar Teman")
		fmt.Println("4. 🤝 Tambah Teman")
		fmt.Println("5. 👤 Profil Saya")
		fmt.Println("6. 🚪 Logout Aplikasi")

		var choice int
		fmt.Print("\n🎯 Pilihan: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			fmt.Println("Beranda (Coming soon)")
		case 2:
			fmt.Println("Apa yang anda pikirkan? Posting sesuatu disini (Coming soon)")
		case 3:
			fmt.Println("Daftar Teman (Coming soon)")
		case 4:
			fmt.Println("Tambah Teman (Coming soon)")
		case 5:
			fmt.Println("Profil Saya (Coming soon)")
		case 6:
			currentUser = nil
			fmt.Println("Logout berhasil")
			return
		}
	}
}

func main() {
	for {
		showHeader("Welcome")
		fmt.Println("1. 🔑 Login")
		fmt.Println("2. ❌ Exit")

		var choice int
		fmt.Print("\n🎯 Pilihan: ")
		fmt.Scan(&choice)

		switch choice {
		case 1:
			if login() {
				showMainMenu()
			}
		case 2:
			fmt.Println("👋 Terima kasih telah menggunakan aplikasi!")
			return
		}
	}
}
