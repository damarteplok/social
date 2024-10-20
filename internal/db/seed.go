package db

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/damarteplok/social/internal/store"
)

var usernames = []string{
	"alice", "bob", "charlie", "david", "eve",
	"frank", "grace", "hank", "irene", "jack",
	"kate", "leo", "mike", "nina", "oscar",
	"paul", "quinn", "rachel", "sam", "tina",
	"ursula", "victor", "will", "xena", "yara",
	"zack", "amber", "brian", "cathy", "daniel",
	"emma", "felix", "gwen", "harry", "isla",
	"james", "kelly", "luke", "mia", "noah",
	"olivia", "peter", "quincy", "rose", "sophie",
	"tom", "uma", "vicky", "wade", "xander",
}

var titles = []string{
	"Tech Trends 2024",
	"Healthy Living Tips",
	"Travel Hacks 101",
	"Budgeting Basics",
	"Quick Dinner Ideas",
	"Fitness Myths Busted",
	"Gardening for All",
	"DIY Home Projects",
	"Mindfulness Matters",
	"Pet Care Essentials",
	"Digital Marketing",
	"Sustainable Living",
	"Personal Finance",
	"Food Photography",
	"Fashion Tips 2024",
	"Mental Health Boost",
	"Mobile Apps Review",
	"Home Office Setup",
	"Cooking Made Easy",
	"Photography Tips",
	"Investing Basics",
	"Life Lessons Learned",
	"Home Decor Ideas",
	"Winter Fashion Tips",
	"Best Travel Dest.",
	"Effective Workouts",
	"Healthy Snack Ideas",
	"Gaming News Today",
	"Crafting for Fun",
	"Eco-Friendly Tips",
	"Business Strategies",
	"Traveling Alone",
	"Parenting Hacks",
	"Startup Insights",
	"Culinary Adventures",
	"Book Recommendations",
	"Weekly Motivation",
	"Coding for Beginners",
	"Art Inspiration",
	"Podcast Reviews",
	"Career Development",
	"Tech Gadgets 2024",
	"Online Learning",
	"Virtual Events Guide",
	"Latest Movie Picks",
	"Self-Care Routines",
	"Cooking Tips",
	"Productivity Hacks",
	"Fitness Challenges",
	"Nature Exploration",
}

var contents = []string{
	"Teknologi terus berkembang dengan cepat. Di tahun 2024, kita akan melihat lebih banyak inovasi dalam AI, perangkat wearable, dan Internet of Things. Pastikan untuk tetap update dengan tren terbaru agar tidak ketinggalan!",
	"Memulai gaya hidup sehat tidak harus sulit. Cobalah untuk menambahkan lebih banyak sayuran ke dalam dietmu, berolahraga secara teratur, dan cukup tidur. Setiap langkah kecil akan membawamu lebih dekat ke kesehatan yang optimal.",
	"Sebelum melakukan perjalanan, pastikan untuk memeriksa aplikasi perjalanan untuk menemukan penawaran terbaik. Selain itu, bawa barang-barang penting di bagasi kabin untuk menghindari kehilangan barang berharga.",
	"Mengelola keuangan bisa jadi tantangan, tetapi dengan anggaran yang baik, kamu bisa mencapai tujuan keuanganmu. Catat semua pengeluaranmu dan alokasikan dana untuk tabungan agar kamu lebih siap menghadapi situasi darurat.",
	"Tidak punya banyak waktu untuk memasak? Cobalah resep pasta cepat dengan saus tomat dan sayuran. Hanya butuh 30 menit, dan kamu bisa menambahkan protein seperti ayam atau tahu untuk meningkatkan nutrisinya.",
	"Banyak mitos seputar kebugaran yang bisa menyesatkan. Salah satunya adalah bahwa kamu harus berolahraga setiap hari untuk mendapatkan hasil. Sebenarnya, istirahat sama pentingnya untuk pemulihan otot.",
	"Berkebun bisa menjadi hobi yang menyenangkan dan bermanfaat. Mulailah dengan tanaman yang mudah dirawat seperti herbal atau sayuran. Kebun kecil di halaman atau bahkan pot di balkon bisa menjadi tempat yang menyegarkan.",
	"Menciptakan proyek DIY di rumah tidak hanya hemat biaya, tetapi juga bisa sangat memuaskan. Cobalah membuat rak dinding atau lukisan sederhana untuk menambahkan sentuhan personal pada dekorasi rumahmu.",
	"Praktik mindfulness dapat membantu mengurangi stres dan meningkatkan fokus. Luangkan beberapa menit setiap hari untuk bernafas dalam-dalam dan hadir dalam momen. Ini dapat meningkatkan kesejahteraan mentalmu secara keseluruhan.",
	"Merawat hewan peliharaan memerlukan komitmen dan pengetahuan. Pastikan untuk memberi makan makanan berkualitas, rutin membawa mereka ke dokter hewan, dan memberikan banyak cinta serta perhatian. Hewan peliharaan yang bahagia akan menjadi teman setia!",
}

var tags = []string{
	"teknologi", "kesehatan", "travel", "keuangan", "masakan",
	"kebugaran", "berkebun", "DIY", "mindfulness", "hewan peliharaan",
	"fashion", "seni", "motivasi", "pendidikan", "game",
	"startup", "kesejahteraan", "pengembangan diri", "desain", "bisnis",
}

var comments = []string{
	"Sangat informatif, terima kasih!",
	"Saya setuju dengan pendapat ini.",
	"Apa pendapatmu tentang topik ini di tahun depan?",
	"Ini adalah tips yang sangat berguna!",
	"Saya suka cara penulis menyampaikan ide.",
	"Berharap ada lebih banyak konten seperti ini.",
	"Bisakah kamu menjelaskan lebih lanjut tentang bagian ini?",
	"Saya sudah mencoba tips ini dan berhasil!",
	"Keren, saya akan mencobanya segera!",
	"Apakah ada saran lain yang bisa kamu bagi?",
	"Bagaimana dengan pengalaman pribadi kamu?",
	"Menarik! Saya ingin tahu lebih banyak.",
	"Ini membuka wawasan saya tentang topik ini.",
	"Bagus sekali! Teruslah berkarya.",
	"Saya sangat menghargai informasi ini.",
	"Apa ada penelitian terbaru tentang hal ini?",
	"Ini sangat membantu, terima kasih!",
	"Saya suka membaca artikel-artikel seperti ini.",
	"Apakah kamu punya buku rekomendasi tentang topik ini?",
	"Komentar ini sangat mencerahkan!",
	"Tunggu, ada yang bisa ditambahkan dari perspektif lain.",
}

func Seed(store store.Storage) {
	ctx := context.Background()

	users := generateUsers(100)
	for _, user := range users {
		if err := store.Users.Create(ctx, user); err != nil {
			log.Println("Error creating user:", err)
			return
		}
	}
	posts := generatePosts(200, users)
	for _, post := range posts {
		if err := store.Posts.Create(ctx, post); err != nil {
			log.Println("Error creating posts:", err)
			return
		}
	}
	comments := generateComments(500, users, posts)
	for _, comment := range comments {
		if err := store.Comments.Create(ctx, comment); err != nil {
			log.Println("Error creating comments:", err)
			return
		}
	}
	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Username: usernames[i%len(usernames)] + fmt.Sprintf("%d", i),
			Email:    usernames[i%len(usernames)] + fmt.Sprintf("%d", i) + "@example.com",
			Password: "123123",
		}
	}

	return users
}

func generatePosts(num int, users []*store.User) []*store.Post {
	posts := make([]*store.Post, num)

	for i := 0; i < num; i++ {
		user := users[rand.Intn(len(users))]

		posts[i] = &store.Post{
			UserID:  user.ID,
			Title:   titles[rand.Intn(len(titles))],
			Content: contents[rand.Intn(len(contents))],
			Tags: []string{
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
				tags[rand.Intn(len(tags))],
			},
		}
	}

	return posts
}

func generateComments(num int, users []*store.User, posts []*store.Post) []*store.Comment {
	cms := make([]*store.Comment, num)

	for i := 0; i < num; i++ {
		cms[i] = &store.Comment{
			PostID:  posts[rand.Intn(len(posts))].ID,
			UserID:  users[rand.Intn(len(users))].ID,
			Content: comments[rand.Intn(len(comments))],
		}
	}

	return cms
}
