package gae

type GaeAuthenicate struct {
    title string
    found bool
    cancelled bool
}

func (* GaeAuthenicate) get_code {

}
    def get_code(self):
        if self.found:
            return self.title[13:]
        else:
            return None