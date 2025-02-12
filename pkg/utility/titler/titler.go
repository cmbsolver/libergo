package titler

import (
	"config"
	"fmt"
)

func PrintTitle(title string) {
	configuration, _ := config.LoadConfig()
	if configuration.HideTitle {
		return
	}

	fmt.Println("                                                                                                              ")
	fmt.Println("                                                                                                              ")
	fmt.Println("                                                                                                              ")
	fmt.Println("   ########        ###        ########            #   %%                                                      ")
	fmt.Println(" %##  # ####%%%%# %###  #%##% ##%  ##%#  %%       #   %%                                                      ")
	fmt.Println("   ##############%## #####%%%##%##%%#    %%       %#  %%#%%%#     #%%%#   %# %%#  #%%%##%    #%%%#            ")
	fmt.Println("      ## ###     #% #       %#####       %%       @#  %@    %@  #@#   %%  %@#    %%    @%  #@#   %@           ")
	fmt.Println("        %###      ####      ##%%         %%       @#  %%    #@# %@@@@@@@  %%     @#    %%  @%     @%          ")
	fmt.Println("          #####   %## ##  #%%#           %%       @#  %%    #@  %%        %%     @#    %%  %%     @#          ")
	fmt.Println("                  ##                     %@#####  @#  %%@##%@#   %@%##%#  %%     #@%##%%%   %@##%@#           ")
	fmt.Println("                   ##                                                                  %%                     ")
	fmt.Println("                                                                                 %####%@#                     ")
	fmt.Println("                                                                                                              ")
	fmt.Println("==============================================================================================================")
	fmt.Println(title)
	fmt.Println("==============================================================================================================")
}
