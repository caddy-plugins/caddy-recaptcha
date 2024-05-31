package recaptcha

const endpoint = `https://www.recaptcha.net/recaptcha/api/siteverify`

// const endpoint=`https://www.google.com/recaptcha/api/siteverify`

const jslib = `https://www.recaptcha.net/recaptcha/api.js?render=`

const jsscript = `<script src="{{jslib}}{{sitekey}}"></script>
<script>
window.addEventListener('load',function(){
	var form=document.querySelector("form");
	if(!form) return;
	form.addEventListener('submit',function(e) {
		e.preventDefault();
		var that=this;
		grecaptcha.ready(function() {
			grecaptcha.execute('{{sitekey}}', {action: '{{action}}'}).then(function(token) {
				var eg=that.querySelector("input[name='g-recaptcha-response']")
				if(!eg){
					var el=document.createElement('input');
					el.type='hidden';
					el.name='g-recaptcha-response';
					el.value=token;
					that.appendChild(el);
				}else{
					eg.value=token;
				}
				that.submit();
			});
		});
	})
})
</script>`
