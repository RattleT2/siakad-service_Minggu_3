def kalkulator ():
    print("Selamat datang di kalkulator sederhana!")
    print("Pilih operasi yang ingin dilakukan:")
    print("1. Penjumlahan")
    print("2. Pengurangan")
    print("3. Perkalian")
    print("4. Pembagian")
    print("5. Keluar")

    pilihan = input("Masukkan pilihan (1/2/3/4/5): ")
    
    while pilihan != '5':
        match pilihan:
            case '1':
                angka1 = float(input("Masukkan angka pertama: "))
                angka2 = float(input("Masukkan angka kedua: "))
                hasil = angka1 + angka2
                print(f"Hasil penjumlahan: {hasil}")
            case '2':
                angka1 = float(input("Masukkan angka pertama: "))
                angka2 = float(input("Masukkan angka kedua: "))
                hasil = angka1 - angka2
                print(f"Hasil pengurangan: {hasil}")
            case '3':
                angka1 = float(input("Masukkan angka pertama: "))
                angka2 = float(input("Masukkan angka kedua: "))
                hasil = angka1 * angka2
                print(f"Hasil perkalian: {hasil}")
            case '4':
                angka1 = float(input("Masukkan angka pertama: "))
                angka2 = float(input("Masukkan angka kedua: "))
                if angka2 != 0:
                    hasil = angka1 / angka2
                    print(f"Hasil pembagian: {hasil}")
                else:
                    print("Error: Pembagian dengan nol tidak diperbolehkan.")
            case _:
                print("Pilihan tidak valid. Silakan coba lagi.")
        
        pilihan = input("Masukkan pilihan (1/2/3/4/5): ")
    
    print("Terima kasih telah menggunakan kalkulator sederhana!")

kalkulator()