# The in, the now, the hip, the with it

- YAML files should represent potential "userland" as much as possible
  - e.g. making profile.yml.tmpl visible / in the new repo is to allow the user to adjust build / deploy processes
  - run commands w/n go as much as possible to keep things testable etc


- Consider moving "code" into it's own top level folder, should be able to use utils fine via the namespace
  - allows for individual test files per code plan, probably a pattern for individual files overall


- Using "ValidResponses" as a suffix for test structures is weird, you're actually tracking arguments sent to it, not the method responses
  - this is currently setup for at least runCmdDir_test and static_test files for now
    - this whole structure is growing rapidly. would be nice to figure out an approach w/ generics if golang has them or templating if not.
      - it's also somewhat questionable. might be patient here for a bit. look into pre existing stubs / spies libs


- A GH Issue exists but much confusion between paths and profiles and the idea of an environment -- need to consider approaches
  - Currently trying to mimic old OMGD rust code in Go, makes sense save for the expansion of these concepts which will expand any refactoring
  - Figuring it out now minimizes changes but may take awhile as I haven't really gotten my head back around this library, had a different concept in mind originally given separation b/t this and OMGD as a rust lib, which, was exploratory, but anyway, it's a cave I need to get out of (amongst others)





# A hopeful future

- Might be soon but next thing post consolidation to think of may be how to split code against different setups / engines e.g. unity et. al. -- keep w/ godot at first, etc
  - Would be really nice to be modular wrt game engines & nakama server but maybe unnecessary, also nakama + agones may be solid no matter the game engine
    - Either way it would be nice as modularity at least enfers less rewriting / copying of server code
  - This does not mean writing the code for other game engines etc right now -- just setting up a folder & naming structure for such things
  - Remember to cater for version #'s -- you're on Godot 3, Godot 4 is a whole other engine, think of semver a bit, etc


- Setting up Agones before approaching other game engines has both risks and rewards
  - The rewards being fun. See if you can setup Agones on a single docker-compose setup -- at least locally. Might be good enough for staging / demo servers, which, would be good enough for quite sometime. That said k8s should be terraform / docker approachable on gcp easy enough.
    - Making a suitable physics based example may be tough


- Where does this stand as a production tool? Not yet I'd assume but it's not far off
  - Very nice to have problem





## OLD

- Using internal go code vs. issuing CLI commands back into the program
  - Using internal go code as much as possible
    - Reinforces tests / structure
    - Compiler will break if internal interfaces are changed, forcing the above to change
    - More work, better for catching further inconsistencies
  - Issuing CLI commands back into the program
    - Potential OS issues
    - Much easier, both to enforce as needed as well as maintain
    - Not exactly unheard of, but seems less organized
  - Decision -- gonna use internal go code as much as possible for as long as maintainable
