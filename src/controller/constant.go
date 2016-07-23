package controller

const (
	firstComment = `<div class="comment-node clear-g block-cont-gw block-cont-bg" datatype="time" cid={{.CommentID}}> 
	<div class="cont-head-gw"> 
		<div class="head-img-gw">
			<img src={{.User.SmallFigureurl}} width="42" height="42" uid="{{.UserID}}"></img>
        </div> 
	</div> 
	<div class="cont-msg-gw"> 
		<div class="msg-wrap-gw"> 
			<div class="wrap-user-gw global-clear-spacing"> 
				<span class="user-time-gw user-time-bg evt-time">{{.CommentTime}}</span> 
				<span class="user-name-gw" title="{{.User.UserName}}"><a href="javascript:void(0)" href="{{.User.SmallFigureurl}}" uid="-1990645212">{{.User.UserName}}</a></span> 
			</div> 
			{{.ChildContent}}
			<div class="wrap-issue-gw"> 
				<p class="issue-wrap-gw"> <span class="wrap-word-bg ">{{.CommentContent}}</span> </p> 
			</div> 
			<div class="clear-g wrap-action-gw"> 
				<div class="action-click-gw"> 
					<i class="gap-gw"></i> 
					<span class="click-ding-gw"><a href="javascript:void(0)" title="顶" class="evt-support"><i class="icon-gw icon-ding-bg"></i><em class="icon-name-bg"></em></a></span>
					<i class="gap-gw"></i>
					<span class="click-cai-gw"><a href="javascript:void(0)" title="踩" class="evt-opposed"><i class="icon-gw icon-cai-bg"></i><em class="icon-name-bg"></em></a></span>
					<i class="gap-gw"></i>
					<span class="click-reply-gw click-reply-eg"><a href="javascript:void(0)" class="evt-reply">回复</a></span>
					<i class="gap-gw"></i>
					<span class="click-share-gw click-reply-eg"><a href="javascript:void(0)" class="evt-share">分享</a></span> 
				</div> 
				<div class="action-from-gw action-from-bg"></div> 
			</div> 
			<div class="wrap-reply-gw"></div> 
			<div class="module-cmt-box" style="display: none;">
			    <!-- 展开状态 -->
			    <div class="post-wrap-w">
			        <div class="wrap-area-w">
			            <div class="area-textarea-w">
			                <textarea node-type="textarea" name="" class="textarea-fw textarea-bf">有事没事说两句...</textarea>
			            </div>
			        </div>
			        <div class="clear-g wrap-action-w">
			            <div class="clear-g action-issue-w">
			                <div class="issue-btn-w">
			                    <a href="javascript:void(0)">
			                        <button node-type="issue" class="btn-fw btn-reply">发布</button>
			                    </a>
			                </div>
			                <div class="issue-icon-w" node-type="share-icons">
			                </div>
			            </div>
			        </div>
			    </div>
			</div>
		</div> 
	</div> 
</div>`
	secondComment = `<div class="wrap-build-gw"> 
	<div class="build-floor-gw"> 
		<div class="build-msg-gw borderbot" cid={{.CommentID}} floornum="1"> 
			{{.ChildContent}}
			<div class="wrap-user-gw global-clear-spacing"> 
				<span class="user-time-gw user-time-bg user-floor-gw">{{.Floor}}</span> 
				<span class="user-name-gw">
					<img style="height: 44px; width: 44px;" src="{{.User.SmallFigureurl}}" title="{{.User.UserName}}" uid="{{.UserID}}">
						<a href="javascript:void(0)" style="margin-left: 5px;">{{.User.UserName}}</a>
					</img>
				</span> 
			</div> 
			<div class="wrap-issue-gw"> 
				<p class="issue-wrap-gw"> <span class="wrap-word-bg ">{{.CommentContent}}</span> </p> 
			</div> 
			<div class="comment-node clear-g wrap-action-gw evt-active-wrapper" style="visibility: hidden;"> 
				<div class="action-click-gw"> 
					<i class="gap-gw"></i> 
					<span class="click-ding-gw"><a href="javascript:void(0)" title="顶" class="evt-support"><i class="icon-gw icon-ding-bg"></i><em class="icon-name-bg"></em></a></span> 
					<i class="gap-gw"></i> 
					<span class="click-cai-gw"><a href="javascript:void(0)" title="踩" class="evt-opposed"><i class="icon-gw icon-cai-bg"></i><em class="icon-name-bg"></em></a></span> 
					<i class="gap-gw"></i> 
					<span class="click-reply-gw click-reply-eg"><a href="javascript:void(0)" class="evt-reply">回复</a></span> 
					<i class="gap-gw"></i> 
					<span class="click-share-gw click-reply-eg"><a href="javascript:void(0)" class="evt-share">分享</a></span> 
				</div> 
				<div class="action-from-gw action-from-bg"></div> 
			</div> 
			<div class="wrap-reply-gw"></div>
			<div class="module-cmt-box" style="display: none;">
			    <!-- 展开状态 -->
			    <div class="post-wrap-w">
			        <div class="wrap-area-w">
			            <div class="area-textarea-w">
			                <textarea node-type="textarea" name="" class="textarea-fw textarea-bf">有事没事说两句...</textarea>
			            </div>
			        </div>
			        <div class="clear-g wrap-action-w">
			            <div class="clear-g action-issue-w">
			                <div class="issue-btn-w">
			                    <a href="javascript:void(0)">
			                        <button node-type="issue" class="btn-fw btn-reply">发布</button>
			                    </a>
			                </div>
			                <div class="issue-icon-w" node-type="share-icons">
			                </div>
			            </div>
			        </div>
			    </div>
			</div>
		</div> 
	</div> 
</div>`
)
