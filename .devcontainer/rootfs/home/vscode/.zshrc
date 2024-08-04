# ZSH Auto Suggestions
export ZSH_AUTOSUGGEST_STRATEGY=(history completion)
source /usr/share/zsh-autosuggestions/zsh-autosuggestions.zsh

autoload -Uz compinit
compinit

# Enable Starship prompt
eval "$(starship init zsh)"

# Install terramaid completion
eval $(terramaid completion zsh)

# Show Terramaid version
terramaid version
