package main

// func TestGetDbCredentials(t *testing.T) {
// 	in, err := ioutil.TempFile("", "")
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer in.Close()

// 	_, err = io.WriteString(in, "root\n"+"root\n")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	_, err = in.Seek(0, os.SEEK_SET)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	got1, got2, _ := getDbCredentials(in)
// 	want1, want2 := "root", "root"
// 	if got1 != want1 || got2 != want2 {
// 		t.Errorf("got1 %q want1 %q\n\t     got2 %q want2 %q", got1, want1, got2, want2)
// 	}
// }
