pub fn get_directory(record_as: String) {
    if cfg!(debug_assertions) {
        // get the directory from the local filesystem
        println!("Recording to directory {} from local filesystem", record_as)
    } else {
        println!("Recording to directory {} from github repo", record_as)
        // get the directory from github.com
    }
}

